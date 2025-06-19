import { get, handleJSONResponse, post, postFormData } from "@sv/apis/common";
import { BYTES_PER, ROOT_FOLDER_ID, ROOT_FOLDER_NAME } from "@sv/utils/consts";
import type {
  FileInfo,
  FolderContent,
  FolderInfo,
  UploadFileInfo,
  UploadFileResult,
} from "./models";

const urlMedia = "media";
const urlFolders = `${urlMedia}/folders`;

// Utility functions
function generateUploadId(): string {
  return `upload_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`;
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
  onProgress?: (progress: number) => void
): Promise<FileInfo> {
  // TODO: Get the chunk size from the backend
  const chunkSize = 5 * BYTES_PER.MB; // 5MB chunks
  const totalChunks = Math.ceil(file.size / chunkSize);
  const uploadId = generateUploadId();
  let finalFileInfo: FileInfo | null = null;

  for (let chunkIndex = 0; chunkIndex < totalChunks; chunkIndex++) {
    const start = chunkIndex * chunkSize;
    const end = Math.min(start + chunkSize, file.size);
    const chunk = file.slice(start, end);
    const isLastChunk = chunkIndex === totalChunks - 1;

    const formData = new FormData();
    formData.append("chunk", chunk);
    formData.append("uploadId", uploadId);
    formData.append("chunkIndex", chunkIndex.toString());
    formData.append("totalChunks", totalChunks.toString());
    formData.append("fileName", file.name);
    formData.append("folderId", folderId);

    if (chunkIndex === 0) {
      formData.append("fileSize", file.size.toString());
      formData.append("mimeType", file.type || "application/octet-stream");
    }

    const res = await postFormData(
      `${urlFolders}/${folderId}/files/chunks`,
      formData
    );

    if (isLastChunk) {
      // Final chunk returns the FileInfo
      finalFileInfo = await handleJSONResponse<FileInfo>(res);
    }

    const progress = ((chunkIndex + 1) / totalChunks) * 100;
    onProgress?.(progress);
  }

  if (!finalFileInfo) {
    throw new Error("Failed to get file info from final chunk");
  }

  return finalFileInfo;
}

export function uploadFiles(
  uFiles: UploadFileInfo[],
  folderId: string,
  onFileProgress?: (id: string, progress: number) => void
): UploadFileResult[] {
  const files = uFiles.map((uFile) => {
    // TODO: Get the limits from the backend
    // Use chunked upload for files larger than 50MB
    const useChunkedUpload = uFile.file.size > 50 * BYTES_PER.MB;

    const file = useChunkedUpload
      ? uploadFileChunked(uFile.file, folderId, (progress) =>
          onFileProgress?.(uFile.id, progress)
        )
      : uploadFile(uFile.file, folderId, (progress) =>
          onFileProgress?.(uFile.id, progress)
        );

    return {
      clientId: uFile.id,
      file,
    };
  });

  return files;
}
