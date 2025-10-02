import { useNavigate } from "@solidjs/router";
import { CLIENT_URLS, FOLDER_CONTENT_TYPES } from "@sv/utils/consts";
import { createSignal, ParentProps, useContext } from "solid-js";
import CTX, { CtxType, SelectedItem } from "./ctx";

const DOUBLE_TAP_DELAY = 500; // ms

export function CtxProvider(props: ParentProps) {
  const navigate = useNavigate();
  const [tapTimer, setTapTimer] = createSignal<number | null>(null);
  const [selectedItem, setSelectedItem] = createSignal<SelectedItem | null>(
    null
  );

  const handleFolderNavigation = (id: string) => {
    clearSelection();
    navigate(`${CLIENT_URLS.DRIVE}/${id}`);
  };

  const clearSelection = () => {
    setSelectedItem(null);
  };

  const handleTap = (
    tappedItem: SelectedItem,
    singleTapAction?: () => void
  ) => {
    if (tapTimer() === null) {
      // First tap, start timer for potential second tap
      const timer = window.setTimeout(() => {
        // If timer completes, it was a single tap

        // Show context menu by selecting the item
        setSelectedItem(tappedItem);

        if (singleTapAction) {
          singleTapAction();
        }
        setTapTimer(null);
      }, DOUBLE_TAP_DELAY);

      setTapTimer(timer);
    } else {
      // Second tap within the delay - it's a double tap
      window.clearTimeout(tapTimer()!);
      setTapTimer(null);

      // Clear any selection on double tap
      clearSelection();

      // Only navigate if it's a folder
      if (tappedItem.type === FOLDER_CONTENT_TYPES.FOLDER) {
        handleFolderNavigation(tappedItem.id);
      }
    }
  };

  const val: CtxType = {
    handleTap,
    selectedItem,
    clearSelection,
  };

  return <CTX.Provider value={val}>{props.children}</CTX.Provider>;
}

function useCtx() {
  const ctx = useContext(CTX);
  if (!ctx) {
    throw new Error("folderContent: useCtx must be used within a CtxProvider");
  }
  return ctx;
}

export default useCtx;
