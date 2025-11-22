import { StorageUsage } from "@sv/apis/profile";
import { SystemConfig } from "@sv/apis/system";
import { Accessor, createContext } from "solid-js";

export const DefaultSystemConfig: SystemConfig = {
  maxDirectUploadSizeMB: 50,
  maxChunkSizeMB: 10,
} as const;

export const DefaultStorageUsage: StorageUsage = {
  used: 0,
  quota: 0,
} as const;

export interface AppCtxType {
  currentFolderId: Accessor<string>;
  systemConfig: SystemConfig;
  storageUsage: Accessor<StorageUsage>;
  refreshStorageUsage: () => void;
}

const AppCtx = createContext<AppCtxType>();

export default AppCtx;
