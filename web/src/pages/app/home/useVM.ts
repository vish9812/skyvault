import { createResource, createSignal } from "solid-js";
import { CONTENT_VIEWS, LOCAL_STORAGE_KEYS } from "@sv/utils/consts";
import { fetchFolderContent } from "@sv/apis/media";
import { useParams } from "@solidjs/router";

const defaultView =
  (localStorage.getItem(
    LOCAL_STORAGE_KEYS.CONTENT_VIEW
  ) as (typeof CONTENT_VIEWS)[keyof typeof CONTENT_VIEWS]) ||
  CONTENT_VIEWS.LIST;

function useVM() {
  const params = useParams();

  const [folderContentRes] = createResource(
    () => params.folderId || "0",
    fetchFolderContent
  );

  const [isListView, setIsListView] = createSignal(
    defaultView === CONTENT_VIEWS.LIST
  );

  const handleListViewChange = () => {
    const isNewViewList = !isListView();
    localStorage.setItem(
      LOCAL_STORAGE_KEYS.CONTENT_VIEW,
      isNewViewList ? CONTENT_VIEWS.LIST : CONTENT_VIEWS.GRID
    );
    setIsListView(isNewViewList);
  };

  return {
    isListView,
    setIsListView,
    handleListViewChange,
    folderContentRes,
  };
}

export default useVM;
