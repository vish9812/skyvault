import { useNavigate } from "@solidjs/router";
import { CLIENT_URLS, FOLDER_CONTENT_TYPES } from "@sv/utils/consts";
import { createSignal } from "solid-js";

const DOUBLE_TAP_DELAY = 500; // ms

function useViewModel() {
  const navigate = useNavigate();
  const [tapTimer, setTapTimer] = createSignal<number | null>(null);

  const handleFolderNavigation = (id: number) => {
    navigate(`${CLIENT_URLS.HOME}${id}`);
  };

  const handleTap = (
    type: string,
    id: number,
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

  return {
    handleTap,
  };
}

export default useViewModel;
