import { ParentProps } from "solid-js";
import { Portal } from "solid-js/web";

interface OverlayProps extends ParentProps {
  handleClick?: () => void;
}

function Overlay(props: OverlayProps) {
  return (
    <Portal>
      <div class="fixed inset-0 z-50 flex-center">
        <div
          class="fixed inset-0 bg-black/50 transition-opacity"
          onClick={props.handleClick}
        />
        {props.children}
      </div>
    </Portal>
  );
}

export default Overlay;
