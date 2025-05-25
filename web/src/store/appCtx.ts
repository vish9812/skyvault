import { Accessor, createContext } from "solid-js";

export interface AppCtxType {
  currentFolderId: Accessor<string>;
}

const AppCtx = createContext<AppCtxType>();

export default AppCtx;
