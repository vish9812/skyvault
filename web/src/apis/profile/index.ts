import { get, handleJSONResponse } from "@sv/apis/common";
import { LOCAL_STORAGE_KEYS } from "@sv/utils/consts";
import type { Profile, StorageUsage } from "./models";

export type { StorageUsage } from "./models";

const urlProfile = "profile";

export function getProfile(): Profile {
  // biome-ignore lint/style/noNonNullAssertion: getProfile is only called when logged in
  const profile = localStorage.getItem(LOCAL_STORAGE_KEYS.PROFILE)!;
  return JSON.parse(profile);
}

export async function fetchStorageUsage(): Promise<StorageUsage> {
  const res = await get(`${urlProfile}/storage`);
  return handleJSONResponse<StorageUsage>(res);
}
