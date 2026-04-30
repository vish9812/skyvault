import type { StorageUsage } from "@sv/apis/profile";
import type { SystemConfig } from "@sv/apis/system";
import { createContext } from "solid-js";

export const DefaultSystemConfig: SystemConfig = {
  maxDirectUploadSizeMB: 50,
  maxChunkSizeMB: 10,
} as const;

export const DefaultStorageUsage: StorageUsage = {
  used: 0,
  quota: 0,
} as const;

export interface AppCtxType {
  currentFolderId: () => string;
  systemConfig: () => SystemConfig;
  storageUsage: () => StorageUsage;
  refreshStorageUsage: () => void;
  folderContentVersion: () => number;
  refreshFolderContent: () => void;
}

const AppCtx = createContext<AppCtxType>();

export default AppCtx;
