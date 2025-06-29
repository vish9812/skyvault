import { SystemConfig } from "@sv/apis/system";
import { Accessor, createContext } from "solid-js";

export const DefaultSystemConfig: SystemConfig = {
  maxUploadSizeMB: 100,
  maxDirectUploadSizeMB: 50,
  maxChunkSizeMB: 10,
} as const;

export interface AppCtxType {
  currentFolderId: Accessor<string>;
  systemConfig: SystemConfig;
}

const AppCtx = createContext<AppCtxType>();

export default AppCtx;
