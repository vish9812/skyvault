import { ROOT_FOLDER_ID, ROOT_FOLDER_NAME } from "@sv/utils/consts";
import type { FileInfo, FolderContent, FolderInfo } from "./models";
import { get, handleJSONResponse, post, postFormData } from "@sv/apis/common";

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
