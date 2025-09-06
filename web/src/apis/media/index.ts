import { get, handleJSONResponse, post, postFormData, ROOT_URL } from "@sv/apis/common";
import { BYTES_PER, ROOT_FOLDER_ID, ROOT_FOLDER_NAME, LOCAL_STORAGE_KEYS } from "@sv/utils/consts";
import Random from "@sv/utils/random";
import type {
  FileInfo,
  FolderContent,
  FolderInfo,
  UploadFileInfo,
  UploadFileResult,
} from "./models";

const urlMedia = "media";
const urlFolders = `${urlMedia}/folders`;

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
  const chunkSize = 2 * BYTES_PER.MB;
  const totalChunks = Math.ceil(file.size / chunkSize);
  const uploadId = Random.id();

  // Create all chunk upload promises in parallel
  const chunkPromises: Promise<Response>[] = [];
  let completedChunks = 0;

  for (let chunkIndex = 0; chunkIndex < totalChunks; chunkIndex++) {
    const start = chunkIndex * chunkSize;
    const end = Math.min(start + chunkSize, file.size);
    const chunk = file.slice(start, end, file.type);

    const formData = new FormData();
    formData.append("chunk", chunk);
    formData.append("uploadId", uploadId);
    formData.append("chunkIndex", chunkIndex.toString());
    formData.append("totalChunks", totalChunks.toString());

    const chunkPromise = new Promise<Response>(async (resolve, reject) => {
      try {
        const res = await postFormData(
          `${urlFolders}/${folderId}/files/chunks`,
          formData
        );

        // Check for successful chunk upload
        if (!res.ok) {
          throw new Error(`Failed to upload chunk ${chunkIndex}`);
        }

        // Update progress for completed chunks
        completedChunks++;
        if (onProgress) {
          // Progress during chunk uploads (90% of total progress)
          const chunkProgress = (completedChunks / totalChunks) * 90;
          onProgress(chunkProgress);
        }

        resolve(res);
      } catch (err) {
        reject(err);
      }
    });

    chunkPromises.push(chunkPromise);
  }

  // Wait for all chunks to complete
  await Promise.all(chunkPromises);

  // All chunks uploaded successfully, now finalize
  const finalizeRes = await post(
    `${urlFolders}/${folderId}/files/chunks/${uploadId}/finalize`,
    {
      fileName: file.name,
      fileSize: file.size,
      mimeType: file.type || "application/octet-stream",
      totalChunks,
    }
  );

  if (onProgress) {
    // Final 10% for finalization
    onProgress(100);
  }

  return handleJSONResponse<FileInfo>(finalizeRes);
}

export function uploadFiles(
  uFiles: UploadFileInfo[],
  folderId: string,
  onFileProgress?: (id: string, progress: number) => void
): UploadFileResult[] {
  const files = uFiles.map((uFile) => {
    // TODO: Get the limits from the backend
    // Use chunked upload for files larger than 50MB
    const useChunkedUpload = uFile.file.size > 2 * BYTES_PER.MB;

    const file = useChunkedUpload
      ? uploadFileChunked(
          uFile.file,
          folderId,
          onFileProgress
            ? (progress) => onFileProgress(uFile.id, progress)
            : undefined
        )
      : uploadFile(
          uFile.file,
          folderId,
          onFileProgress
            ? (progress) => onFileProgress(uFile.id, progress)
            : undefined
        );

    return {
      clientId: uFile.id,
      file,
    };
  });

  return files;
}

export async function downloadFile(fileId: string, fileName: string): Promise<void> {
  const token = localStorage.getItem(LOCAL_STORAGE_KEYS.TOKEN);
  if (!token) throw new Error("No token found");

  try {
    // Fetch the file with authorization header
    const response = await fetch(`${ROOT_URL}/${urlMedia}/files/${fileId}/download`, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });

    if (!response.ok) {
      throw new Error(`Download failed: ${response.statusText}`);
    }

    // Get the blob from the response
    const blob = await response.blob();
    
    // Create a link element to trigger the download
    const link = document.createElement('a');
    const url = window.URL.createObjectURL(blob);
    
    link.href = url;
    link.download = fileName;
    link.setAttribute('style', 'display: none');
    
    document.body.appendChild(link);
    link.click();
    
    // Clean up
    document.body.removeChild(link);
    window.URL.revokeObjectURL(url);
  } catch (error) {
    console.error('Download failed:', error);
    throw error;
  }
}
