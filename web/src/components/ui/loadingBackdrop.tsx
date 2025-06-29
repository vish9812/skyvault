import Overlay from "./overlay";

interface LoadingBackdropProps {
  message?: string;
}

function LoadingBackdrop(props: LoadingBackdropProps) {
  return (
    <Overlay>
      {/* Spinner */}
      <div class="w-8 h-8 border-4 border-primary-dark rounded-full border-t-primary animate-spin"></div>

      {/* Message */}
      <div class="text-center pl-4">
        <h3 class="text-neutral font-medium">
          {props.message || "Loading..."}
        </h3>
      </div>
    </Overlay>
  );
}

export default LoadingBackdrop;
