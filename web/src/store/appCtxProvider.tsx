import { useLocation, useMatch, useNavigate, useParams } from "@solidjs/router";
import { isLoggedIn } from "@sv/apis/auth";
import { fetchStorageUsage } from "@sv/apis/profile";
import { getSystemConfig } from "@sv/apis/system";
import LoadingBackdrop from "@sv/components/ui/loadingBackdrop";
import { CLIENT_URLS, ROOT_FOLDER_ID } from "@sv/utils/consts";
import {
  createRenderEffect,
  createResource,
  createSignal,
  type ParentProps,
  Show,
  useContext,
} from "solid-js";
import AppCtx, { DefaultStorageUsage, DefaultSystemConfig } from "./appCtx";

export function AppCtxProvider(props: ParentProps) {
  const navigate = useNavigate();

  if (!isLoggedIn()) {
    navigate(CLIENT_URLS.SIGN_IN, { replace: true });
    return;
  }

  // System config only loads when user is authenticated
  const [systemConfig] = createResource(getSystemConfig, {
    initialValue: DefaultSystemConfig,
  });

  // Storage usage
  const [storageUsage, { refetch: refreshStorageUsage }] = createResource(
    () => fetchStorageUsage(),
    {
      initialValue: DefaultStorageUsage,
    },
  );

  // Current folder id
  const location = useLocation();

  const isNavigatable = useMatch(
    () => location.pathname,
    [CLIENT_URLS.DRIVE, CLIENT_URLS.SHARED],
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

  const [folderContentVersion, setFolderContentVersion] = createSignal(0);
  const refreshFolderContent = () => setFolderContentVersion((v) => v + 1);

  return (
    <div>
      <Show
        when={!systemConfig.loading && !storageUsage.loading}
        fallback={<LoadingBackdrop />}
      >
        <AppCtx.Provider
          value={{
            currentFolderId,
            systemConfig,
            storageUsage,
            refreshStorageUsage,
            folderContentVersion,
            refreshFolderContent,
          }}
        >
          {props.children}
        </AppCtx.Provider>
      </Show>
    </div>
  );
}

function useAppCtx() {
  const ctx = useContext(AppCtx);
  if (!ctx) {
    throw new Error("app: useAppCtx must be used within an AppCtxProvider");
  }
  return ctx;
}

export default useAppCtx;
