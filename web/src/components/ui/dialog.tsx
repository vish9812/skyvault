import { createEffect, JSX, Show } from "solid-js";
import Overlay from "./overlay";

interface DialogProps {
  open: boolean;
  onClose: () => void;
  title: string;
  description?: string;
  children: JSX.Element;
  size?: "sm" | "md" | "lg" | "xl";
  actions?: JSX.Element;
}

export default function Dialog(props: DialogProps) {
  const maxWidthClass = () => {
    switch (props.size) {
      case "sm":
        return "max-w-sm";
      case "md":
        return "max-w-md";
      case "lg":
        return "max-w-lg";
      case "xl":
        return "max-w-xl";
      default:
        return "max-w-2xl";
    }
  };

  // Handle escape key to close dialog
  createEffect(() => {
    if (props.open) {
      const handleKeyDown = (e: KeyboardEvent) => {
        if (e.key === "Escape") {
          props.onClose();
        }
      };
      document.addEventListener("keydown", handleKeyDown);

      // Prevent body scroll when dialog is open
      document.body.style.overflow = "hidden";

      return () => {
        document.removeEventListener("keydown", handleKeyDown);
        document.body.style.overflow = "unset";
      };
    }
  });

  return (
    <Show when={props.open}>
      <Overlay handleClick={undefined}>
        {/* Modal Content */}
        {/* TODO: Tailwindcss v4 is missing the animate-in, fade-in-0, and zoom-in-95 classes */}
        <div
          class={`relative z-10 w-full ${maxWidthClass()} mx-4 bg-white rounded-lg shadow-xl max-h-[90vh] overflow-hidden animate-in fade-in-0 zoom-in-95 duration-200`}
          onClick={(e) => e.stopPropagation()}
          onTouchStart={(e) => e.stopPropagation()}
        >
          <div class="flex flex-col">
            {/* Header */}
            <div class="px-6 py-4 border-b border-border-strong">
              <div class="flex items-center justify-between">
                <div>
                  <h2 class="text-xl font-semibold text-neutral">
                    {props.title}
                  </h2>
                  <Show when={props.description}>
                    <p class="mt-1 text-sm text-neutral-lighter">
                      {props.description}
                    </p>
                  </Show>
                </div>
                <button
                  onClick={props.onClose}
                  class="text-neutral-light hover:text-neutral transition-colors"
                >
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke-width="1.5"
                    stroke="currentColor"
                    class="w-6 h-6"
                  >
                    <path
                      stroke-linecap="round"
                      stroke-linejoin="round"
                      d="M6 18 18 6M6 6l12 12"
                    />
                  </svg>
                </button>
              </div>
            </div>

            {/* Content */}
            <div class="px-6 py-4 overflow-y-auto">{props.children}</div>

            {/* Actions */}
            {props.actions && (
              <div class="px-6 py-4 border-t border-border-strong flex justify-end gap-2">
                {props.actions}
              </div>
            )}
          </div>
        </div>
      </Overlay>
    </Show>
  );
}
