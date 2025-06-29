import { createSignal, ParentProps, useContext } from "solid-js";
import CTX, { CtxType } from "./ctx";
import { CLIENT_URLS, FOLDER_CONTENT_TYPES } from "@sv/utils/consts";
import { useNavigate } from "@solidjs/router";

const DOUBLE_TAP_DELAY = 500; // ms

export function CtxProvider(props: ParentProps) {
  const navigate = useNavigate();
  const [tapTimer, setTapTimer] = createSignal<number | null>(null);

  const handleFolderNavigation = (id: string) => {
    navigate(`${CLIENT_URLS.DRIVE}/${id}`);
  };

  const handleTap = (
    type: string,
    id: string,
    singleTapAction?: () => void
  ) => {
    if (tapTimer() === null) {
      // First tap, start timer for potential second tap
      const timer = window.setTimeout(() => {
        // If timer completes, it was a single tap
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

      // Only navigate if it's a folder
      if (type === FOLDER_CONTENT_TYPES.FOLDER) {
        handleFolderNavigation(id);
      }
    }
  };

  const val: CtxType = {
    handleTap,
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
