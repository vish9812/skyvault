import type { FolderContent, FolderInfo } from "./models";
import { get, handleJSONResponse, post } from "@sv/apis/common";

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
