import {
  del,
  get,
  handleBlobResponse,
  handleJSONResponse,
  patch,
  post,
  postFormData,
} from "@sv/apis/common";
import { BYTES_PER, ROOT_FOLDER_ID, ROOT_FOLDER_NAME } from "@sv/utils/consts";
import FileUtils from "@sv/utils/fileUtils";
import type {
  CreateUploadSessionResponse,
  FileInfo,
  FolderContent,
  FolderInfo,
  Page,
  UploadConfig,
  UploadFileInfo,
  UploadFileResult,
} from "./models";

const urlMedia = "media";
const urlFolders = `${urlMedia}/folders`;
const urlFiles = `${urlMedia}/files`;

interface PagingParams {
  limit?: number;
  direction?: string;
  sort?: string;
  sortBy?: string;
  nextCursor?: string;
  prevCursor?: string;
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

export async function fetchChildFolders(
  id: string,
  paging?: PagingParams
): Promise<Page<FolderInfo>> {
  const res = await get(`${urlFolders}/${id}/folders${queryString(paging)}`);
  return handleJSONResponse<Page<FolderInfo>>(res);
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
  uploadConfig: UploadConfig,
  file: File,
  folderId: string,
  onProgress?: (progress: number) => void
): Promise<FileInfo> {
  const maxChunkSize = uploadConfig.maxChunkSizeMB * BYTES_PER.MB;
  const totalChunks = Math.ceil(file.size / maxChunkSize);

  // Step 1: Create upload session (allocates quota upfront)
  const sessionRes = await post(
    `${urlFolders}/${folderId}/files/upload-sessions`,
    {
      fileName: file.name,
      fileSize: file.size,
      mimeType: file.type || "application/octet-stream",
      totalChunks,
    }
  );
  const session =
    await handleJSONResponse<CreateUploadSessionResponse>(sessionRes);
  const uploadId = session.uploadId;

  const chunksUrl = `${urlFolders}/${folderId}/files/chunks`;

  // Step 2: Create all chunk upload promises in parallel
  const chunkPromises: Promise<Response>[] = [];
  let completedChunks = 0;

  for (let chunkIndex = 0; chunkIndex < totalChunks; chunkIndex++) {
    const start = chunkIndex * maxChunkSize;
    const end = Math.min(start + maxChunkSize, file.size);
    const chunk = file.slice(start, end, file.type);

    const formData = new FormData();
    formData.append("chunk", chunk);
    formData.append("uploadId", uploadId);
    formData.append("chunkIndex", chunkIndex.toString());
    formData.append("totalChunks", totalChunks.toString());

    const chunkPromise = new Promise<Response>(async (resolve, reject) => {
      try {
        const res = await postFormData(chunksUrl, formData);

        // Check for successful chunk upload
        if (!res.ok) {
          throw new Error(`Failed to upload chunk ${chunkIndex}`);
        }

        // Update progress for completed chunks
        completedChunks++;
        if (onProgress) {
          // Chunk uploads represent 90% of progress; finalization is the remaining 10% which happens further down the code
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

  // Step 3: All chunks uploaded successfully, now finalize
  const finalizeRes = await post(`${chunksUrl}/${uploadId}/finalize`, {
    fileName: file.name,
    fileSize: file.size,
    mimeType: file.type || "application/octet-stream",
    totalChunks,
  });

  if (onProgress) {
    // Final 10% for finalization
    onProgress(100);
  }

  return handleJSONResponse<FileInfo>(finalizeRes);
}

export function uploadFiles(
  uploadConfig: UploadConfig,
  uFiles: UploadFileInfo[],
  folderId: string,
  onFileProgress?: (id: string, progress: number) => void
): UploadFileResult[] {
  const files = uFiles.map((uFile) => {
    const useChunkedUpload =
      uFile.file.size > uploadConfig.maxDirectUploadSizeMB * BYTES_PER.MB;

    const file = useChunkedUpload
      ? uploadFileChunked(
          uploadConfig,
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

export async function downloadFile(
  fileId: string,
  fileName: string
): Promise<void> {
  const response = await post(`${urlFiles}/${fileId}/download`, {});
  const blob = await handleBlobResponse(response);
  FileUtils.downloadBlob(blob, fileName);
}

export async function renameFile(fileId: string, name: string): Promise<void> {
  const response = await patch(`${urlFiles}/${fileId}/rename`, { name });
  await handleEmptyResponse(response);
}

export async function renameFolder(
  folderId: string,
  name: string
): Promise<void> {
  const response = await patch(`${urlFolders}/${folderId}/rename`, { name });
  await handleEmptyResponse(response);
}

export async function moveFile(
  fileId: string,
  folderId: string
): Promise<void> {
  const response = await patch(`${urlFiles}/${fileId}/move`, { folderId });
  await handleEmptyResponse(response);
}

export async function moveFolder(
  folderId: string,
  destinationFolderId: string
): Promise<void> {
  const response = await patch(`${urlFolders}/${folderId}/move`, {
    folderId: destinationFolderId,
  });
  await handleEmptyResponse(response);
}

export async function trashFiles(fileIds: string[]): Promise<void> {
  const response = await del(`${urlFiles}/`, { fileIds });
  await handleEmptyResponse(response);
}

export async function trashFolders(folderIds: string[]): Promise<void> {
  const response = await del(`${urlFolders}/`, { folderIds });
  await handleEmptyResponse(response);
}

async function handleEmptyResponse(res: Response): Promise<void> {
  if (res.status === 204) return;
  await handleJSONResponse<unknown>(res);
}

function queryString(params?: PagingParams): string {
  if (!params) return "";

  const query = new URLSearchParams();
  if (params.limit !== undefined) query.set("limit", String(params.limit));
  if (params.direction) query.set("direction", params.direction);
  if (params.sort) query.set("sort", params.sort);
  if (params.sortBy) query.set("sort-by", params.sortBy);
  if (params.nextCursor) query.set("next-cursor", params.nextCursor);
  if (params.prevCursor) query.set("prev-cursor", params.prevCursor);

  const value = query.toString();
  return value ? `?${value}` : "";
}
