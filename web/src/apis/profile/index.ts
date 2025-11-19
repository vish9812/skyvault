import type { StorageUsage } from "./models";
import { get, handleJSONResponse } from "@sv/apis/common";

export async function fetchStorageUsage(
  profileId: string
): Promise<StorageUsage> {
  const res = await get(`profile/${profileId}/storage`);
  return handleJSONResponse<StorageUsage>(res);
}
