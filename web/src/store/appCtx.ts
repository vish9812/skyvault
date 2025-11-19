import { SystemConfig } from "@sv/apis/system";
import { StorageUsage } from "@sv/apis/profile";
import { Accessor, createContext } from "solid-js";

export const DefaultSystemConfig: SystemConfig = {
  maxDirectUploadSizeMB: 50,
  maxChunkSizeMB: 10,
} as const;

export const DefaultStorageUsage: StorageUsage = {
  usedBytes: 0,
  quotaBytes: 0,
  usedMB: 0,
  quotaMB: 0,
} as const;

export interface AppCtxType {
  currentFolderId: Accessor<string>;
  systemConfig: SystemConfig;
  storageUsage: Accessor<StorageUsage>;
  refreshStorageUsage: () => void;
}

const AppCtx = createContext<AppCtxType>();

export default AppCtx;
