import { useLocation, useMatch, useNavigate, useParams } from "@solidjs/router";
import { getProfile } from "@sv/apis/auth";
import { getSystemConfig } from "@sv/apis/system";
import { fetchStorageUsage } from "@sv/apis/profile";
import LoadingBackdrop from "@sv/components/ui/loadingBackdrop";
import { CLIENT_URLS, ROOT_FOLDER_ID } from "@sv/utils/consts";
import {
  createRenderEffect,
  createResource,
  createSignal,
  ParentProps,
  Show,
  useContext,
} from "solid-js";
import AppCtx, { DefaultSystemConfig, DefaultStorageUsage } from "./appCtx";

export function AppCtxProvider(props: ParentProps) {
  const navigate = useNavigate();
  const profile = getProfile();

  if (!profile) {
    navigate(CLIENT_URLS.SIGN_IN, { replace: true });
    return;
  }

  // System config only loads when user is authenticated
  const [systemConfig] = createResource(getSystemConfig, {
    initialValue: DefaultSystemConfig,
  });

  // Storage usage
  const [storageUsage, { refetch: refetchStorageUsage }] = createResource(
    () => fetchStorageUsage(profile.id),
    {
      initialValue: DefaultStorageUsage,
    }
  );

  // Current folder id
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

  return (
    <div>
      <Show
        when={!systemConfig.loading && !storageUsage.loading}
        fallback={<LoadingBackdrop />}
      >
        <AppCtx.Provider
          value={{
            currentFolderId,
            systemConfig: systemConfig()!,
            storageUsage,
            refreshStorageUsage: refetchStorageUsage,
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
