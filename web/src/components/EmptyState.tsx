import { Button } from "@kobalte/core/button";
import type { JSX } from "solid-js";

interface EmptyStateProps {
  title?: string;
  description?: string;
  buttonLabel?: string;
  onButtonClick?: () => void;
}

const EmptyState = (props: EmptyStateProps): JSX.Element => {
  return (
    <div class="text-center py-16 bg-white rounded-lg border border-gray-200 shadow-sm">
      <div class="inline-flex justify-center items-center w-16 h-16 rounded-full bg-primary/20 text-primary mb-4">
        <span class="material-symbols-outlined text-3xl">cloud_upload</span>
      </div>
      <h3 class="text-lg font-medium text-gray-900">
        {props.title || "No files or folders yet"}
      </h3>
      <p class="mt-2 text-sm text-gray-500 max-w-md mx-auto">
        {props.description ||
          "Upload files or create folders to get started with your secure cloud storage."}
      </p>
      <div class="mt-6">
        <Button class="btn btn-primary" onClick={props.onButtonClick}>
          {props.buttonLabel || "Upload Files"}
        </Button>
      </div>
    </div>
  );
};

export default EmptyState;
