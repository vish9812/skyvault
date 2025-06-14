import { useLocation, useMatch, useParams } from "@solidjs/router";
import { CLIENT_URLS, ROOT_FOLDER_ID } from "@sv/utils/consts";
import {
  createRenderEffect,
  createSignal,
  ParentProps,
  useContext
} from "solid-js";
import AppCtx, { AppCtxType } from "./appCtx";

export function AppCtxProvider(props: ParentProps) {
  const location = useLocation();

  const isNavigatable = useMatch(
    () => location.pathname,
    [CLIENT_URLS.DRIVE, CLIENT_URLS.SHARED]
  );

  const params = useParams();
  const [currentFolderId, setCurrentFolderId] = createSignal(ROOT_FOLDER_ID);

  createRenderEffect(() => {
    if (isNavigatable()) {
      setCurrentFolderId(params.folderId || ROOT_FOLDER_ID);
    } else {
      setCurrentFolderId(ROOT_FOLDER_ID);
    }
  });

  const val: AppCtxType = {
    currentFolderId,
  };

  return <AppCtx.Provider value={val}>{props.children}</AppCtx.Provider>;
}

function useAppCtx() {
  const ctx = useContext(AppCtx);
  if (!ctx) {
    throw new Error("app: useAppCtx must be used within an AppCtxProvider");
  }
  return ctx;
}

export default useAppCtx;
