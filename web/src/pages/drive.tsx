import { Button } from "@kobalte/core/button";
import { getProfile } from "@sv/apis/auth";
import { fetchFolderContent, fetchFolderInfo } from "@sv/apis/media";
import Breadcrumbs from "@sv/components/breadcrumbs";
import FolderContent from "@sv/components/folderContent";
import Icon from "@sv/components/icons";
import useAppCtx from "@sv/store/appCtxProvider";
import { Show, createResource, createSignal } from "solid-js";

function Drive() {
  const appCtx = useAppCtx();

  // TODO: Replace with tanstack query
  const [folderContent] = createResource(
    () => appCtx.currentFolderId(),
    fetchFolderContent
  );

  const [folderInfo] = createResource(
    () => appCtx.currentFolderId(),
    fetchFolderInfo
  );

  const [isListView, setIsListView] = createSignal(
    getProfile()!.preferences.contentView === "list"
  );

  const handleContentViewChange = () => {
    const newView = !isListView();

    // TODO: API call to update preference contentView
    // await makeApiCall();

    setIsListView(newView);
  };

  return (
    <>
      <div class="flex items-center justify-between">
        <Show when={!!folderInfo.latest}>
          <Breadcrumbs
            ancestors={folderInfo.latest!.ancestors}
            currentFolder={folderInfo.latest!}
          />
        </Show>
        <div>
          <Button class="btn btn-ghost" onClick={handleContentViewChange}>
            <Show when={isListView()} fallback={<Icon name="list" size={6} />}>
              <Icon name="grid" size={6} />
            </Show>
          </Button>
        </div>
      </div>

      {/* Files & Folders */}
      <FolderContent
        loading={folderContent.loading}
        content={folderContent.latest}
        isListView={isListView()}
      />
    </>
  );
}

export default Drive;
