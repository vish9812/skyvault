import type { FolderContent, FolderInfo } from "./models";
import { getJSON, handleJSONResponse, postJSON } from "@sv/apis/common";

const urlMedia = "media";
const urlFolders = `${urlMedia}/folders`;

export async function fetchFolderContent(
  folderId?: string
): Promise<FolderContent> {
  const id = folderId || "0";

  const res = await getJSON(`${urlFolders}/${id}/content`);

  // Set a timeout to simulate a slow response
  return new Promise((resolve) => {
    setTimeout(async () => {
      const data = await handleJSONResponse(res);
      resolve(data);
    }, 700);
  });
}

export async function createFolder(
  parentFolderId: string | undefined,
  name: string
): Promise<FolderInfo> {
  const id = parentFolderId || "0";
  const res = await postJSON(`${urlFolders}/${id}/`, { name });
  return handleJSONResponse(res);
}
