import { get, handleJSONResponse } from "@sv/apis/common";
import type { StorageUsage } from "./models";

const urlProfile = "profile";

export async function fetchStorageUsage(): Promise<StorageUsage> {
  const res = await get(`${urlProfile}/storage`);
  return handleJSONResponse<StorageUsage>(res);
}
