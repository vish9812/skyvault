import type { FileInfo, FolderContent, FolderInfo } from "./models";
import { get, handleJSONResponse, post, postFormData } from "@sv/apis/common";

const urlMedia = "media";
const urlFolders = `${urlMedia}/folders`;

export async function fetchFolderContent(
  folderId?: string
): Promise<FolderContent> {
  const id = folderId || "0";
  const res = await get(`${urlFolders}/${id}/content`);
  return handleJSONResponse<FolderContent>(res);
}

export async function createFolder(
  parentFolderId: string | undefined,
  name: string
): Promise<FolderInfo> {
  const id = parentFolderId || "0";
  const res = await post(`${urlFolders}/${id}/`, { name });
  return handleJSONResponse<FolderInfo>(res);
}

export async function uploadFile(
  file: File,
  folderId: string | undefined,
  onProgress?: (progress: number) => void
): Promise<FileInfo> {
  const id = folderId || "0";
  const formData = new FormData();
  formData.append("file", file);

  const res = await postFormData(
    `${urlFolders}/${id}/files`,
    formData,
    onProgress
  );
  return handleJSONResponse<FileInfo>(res);
}

export async function uploadFiles(
  files: File[],
  folderId: string | undefined,
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
