import { ROOT_FOLDER_ID, ROOT_FOLDER_NAME } from "@sv/utils/consts";
import type { FileInfo, FolderContent, FolderInfo } from "./models";
import { get, handleJSONResponse, post, postFormData } from "@sv/apis/common";

const urlMedia = "media";
const urlFolders = `${urlMedia}/folders`;

// Utility functions
function generateUploadId(): string {
  return `upload_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`;
}

class Semaphore {
  private permits: number;
  private waiting: (() => void)[] = [];

  constructor(permits: number) {
    this.permits = permits;
  }

  async acquire(): Promise<void> {
    if (this.permits > 0) {
      this.permits--;
      return;
    }

    return new Promise((resolve) => {
      this.waiting.push(resolve);
    });
  }

  release(): void {
    if (this.waiting.length > 0) {
      const resolve = this.waiting.shift()!;
      resolve();
    } else {
      this.permits++;
    }
  }
}

export async function fetchFolderInfo(id: string): Promise<FolderInfo> {
  if (id === ROOT_FOLDER_ID) {
    return Promise.resolve({
      id: ROOT_FOLDER_ID,
      name: ROOT_FOLDER_NAME,
      ownerId: "",
      ancestors: [],
      createdAt: "",
      updatedAt: "",
    });
  }

  const res = await get(`${urlFolders}/${id}`);
  return handleJSONResponse<FolderInfo>(res);
}

export async function fetchFolderContent(id: string): Promise<FolderContent> {
  const res = await get(`${urlFolders}/${id}/content`);
  return handleJSONResponse<FolderContent>(res);
}

export async function createFolder(
  parentFolderId: string,
  name: string
): Promise<FolderInfo> {
  const res = await post(`${urlFolders}/${parentFolderId}/`, { name });
  return handleJSONResponse<FolderInfo>(res);
}

export async function uploadFile(
  file: File,
  folderId: string,
  onProgress?: (progress: number) => void
): Promise<FileInfo> {
  const formData = new FormData();
  formData.append("file", file);

  const res = await postFormData(
    `${urlFolders}/${folderId}/files`,
    formData,
    onProgress
  );
  return handleJSONResponse<FileInfo>(res);
}

export async function uploadFileChunked(
  file: File,
  folderId: string,
  onProgress?: (progress: number) => void,
  chunkSize: number = 5 * 1024 * 1024 // 5MB chunks
): Promise<FileInfo> {
  const totalChunks = Math.ceil(file.size / chunkSize);
  const uploadId = generateUploadId();

  for (let chunkIndex = 0; chunkIndex < totalChunks; chunkIndex++) {
    const start = chunkIndex * chunkSize;
    const end = Math.min(start + chunkSize, file.size);
    const chunk = file.slice(start, end);

    const formData = new FormData();
    formData.append("chunk", chunk);
    formData.append("uploadId", uploadId);
    formData.append("chunkIndex", chunkIndex.toString());
    formData.append("totalChunks", totalChunks.toString());
    formData.append("fileName", file.name);

    if (chunkIndex === 0) {
      formData.append("fileSize", file.size.toString());
      formData.append("mimeType", file.type || "application/octet-stream");
    }

    await postFormData(`${urlFolders}/${folderId}/files/chunks`, formData);

    const progress = ((chunkIndex + 1) / totalChunks) * 100;
    onProgress?.(progress);
  }

  // Finalize upload
  const res = await post(`${urlFolders}/${folderId}/files/finalize`, {
    uploadId,
    fileName: file.name,
    fileSize: file.size,
    mimeType: file.type || "application/octet-stream",
  });

  return handleJSONResponse<FileInfo>(res);
}

export async function uploadFiles(
  files: File[],
  folderId: string,
  onFileProgress?: (fileIndex: number, progress: number) => void
): Promise<FileInfo[]> {
  const results: FileInfo[] = [];

  for (let i = 0; i < files.length; i++) {
    const file = files[i];
    const fileInfo = await uploadFile(file, folderId, (progress) =>
      onFileProgress?.(i, progress)
    );
    results.push(fileInfo);
  }

  return results;
}

export async function uploadFilesParallel(
  files: File[],
  folderId: string,
  onFileProgress?: (fileIndex: number, progress: number) => void,
  maxConcurrent: number = 5 // Limit concurrent uploads
): Promise<FileInfo[]> {
  const results: FileInfo[] = new Array(files.length);
  const semaphore = new Semaphore(maxConcurrent);

  const uploadPromises = files.map(async (file, index) => {
    await semaphore.acquire();
    try {
      // Use chunked upload for files larger than 100MB
      const useChunkedUpload = file.size > 100 * 1024 * 1024;

      const fileInfo = useChunkedUpload
        ? await uploadFileChunked(file, folderId, (progress) =>
            onFileProgress?.(index, progress)
          )
        : await uploadFile(file, folderId, (progress) =>
            onFileProgress?.(index, progress)
          );

      results[index] = fileInfo;
    } finally {
      semaphore.release();
    }
  });

  await Promise.all(uploadPromises);
  return results;
}
