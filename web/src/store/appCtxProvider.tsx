import { useLocation, useMatch, useParams } from "@solidjs/router";
import { CLIENT_URLS, ROOT_FOLDER_ID } from "@sv/utils/consts";
import {
  createEffect,
  createRenderEffect,
  createSignal,
  ParentProps,
  useContext,
} from "solid-js";
import AppCtx, { AppCtxType } from "./appCtx";

export function AppCtxProvider(props: ParentProps) {
  const location = useLocation();
  console.log("location.pathname: ", location.pathname);
  const isNavigatable = useMatch(
    () => location.pathname,
    [CLIENT_URLS.DRIVE, CLIENT_URLS.SHARED]
  );
  console.log("isNavigatable: ", isNavigatable());
  const params = useParams();
  const [currentFolderId, setCurrentFolderId] = createSignal(ROOT_FOLDER_ID);

  createRenderEffect(() => {
    console.log("effect");
    if (isNavigatable()) {
      console.log("params.folderId", params.folderId);
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
