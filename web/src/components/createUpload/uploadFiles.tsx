import { Dialog } from "@kobalte/core/dialog";
import { Meter } from "@kobalte/core/meter";
import { createSignal, For, Show } from "solid-js";
import useCtx from "./ctxProvider";

export default function UploadFiles() {
  const ctx = useCtx();
  const [isDragging, setIsDragging] = createSignal(false);

  const formatFileSize = (bytes: number): string => {
    if (bytes < 1024) return `${bytes} B`;
    if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`;
    if (bytes < 1024 * 1024 * 1024)
      return `${(bytes / (1024 * 1024)).toFixed(1)} MB`;
    return `${(bytes / (1024 * 1024 * 1024)).toFixed(1)} GB`;
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case "success":
        return "text-success";
      case "error":
        return "text-error";
      default:
        return "text-neutral";
    }
  };

  const handleDragOver = (e: DragEvent) => {
    e.preventDefault();
    setIsDragging(true);
  };

  const handleDragLeave = (e: DragEvent) => {
    e.preventDefault();
    setIsDragging(false);
  };

  const handleDrop = (e: DragEvent) => {
    e.preventDefault();
    setIsDragging(false);

    if (e.dataTransfer?.files) {
      ctx.handleFileSelect(e.dataTransfer.files);
    }
  };

  return (
    <Dialog
      open={ctx.isFileUploadModalOpen()}
      onOpenChange={ctx.setIsFileUploadModalOpen}
    >
      <Dialog.Portal>
        <Dialog.Overlay class="dialog-overlay" />
        <Dialog.Content class="dialog-content max-w-lg">
          <div class="flex justify-between items-center mb-4">
            <Dialog.Title class="dialog-title">Upload Files</Dialog.Title>
            <Dialog.CloseButton class="p-1 rounded hover:bg-bg-muted">
              <svg
                xmlns="http://www.w3.org/2000/svg"
                fill="none"
                viewBox="0 0 24 24"
                stroke-width="1.5"
                stroke="currentColor"
                class="size-5"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  d="M6 18 18 6M6 6l12 12"
                />
              </svg>
            </Dialog.CloseButton>
          </div>

          <Show when={ctx.uploadError()}>
            <div class="mb-4 py-2 px-3 bg-error-light/10 text-error rounded">
              {ctx.uploadError()}
            </div>
          </Show>

          <Show
            when={ctx.fileUploads().length > 0}
            fallback={
              <div
                class={`border-2 border-dashed rounded-lg p-6 mb-4 text-center transition-colors ${
                  isDragging()
                    ? "border-primary bg-primary-lighter"
                    : "border-border-strong"
                }`}
                onDragOver={handleDragOver}
                onDragLeave={handleDragLeave}
                onDrop={handleDrop}
              >
                <div class="flex-center flex-col gap-2">
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke-width="1.5"
                    stroke="currentColor"
                    class="size-10 text-neutral-light"
                  >
                    <path
                      stroke-linecap="round"
                      stroke-linejoin="round"
                      d="M3 16.5v2.25A2.25 2.25 0 0 0 5.25 21h13.5A2.25 2.25 0 0 0 21 18.75V16.5m-13.5-9L12 3m0 0 4.5 4.5M12 3v13.5"
                    />
                  </svg>
                  <p class="font-medium">Drag and drop files here</p>
                  <p class="text-sm text-neutral-light">
                    Or click the upload button below
                  </p>
                </div>
              </div>
            }
          >
            <div class="max-h-[300px] overflow-y-auto mb-4">
              <For each={ctx.fileUploads()}>
                {(item) => (
                  <div class="mb-4 pb-4 border-b border-border last:border-b-0 last:pb-0">
                    <div class="flex justify-between mb-1">
                      <div class="truncate max-w-[70%]" title={item.file.name}>
                        {item.file.name}
                      </div>
                      <div class={`text-sm ${getStatusColor(item.status)}`}>
                        {item.status === "pending" && "Pending"}
                        {item.status === "uploading" && `${item.progress}%`}
                        {item.status === "success" && "Completed"}
                        {item.status === "error" && "Failed"}
                      </div>
                    </div>

                    <div class="flex justify-between mb-1 text-sm text-neutral-light">
                      <span>{formatFileSize(item.file.size)}</span>
                      <span>
                        {formatFileSize(item.file.size * (item.progress / 100))}
                      </span>
                    </div>

                    <Meter
                      value={item.progress}
                      minValue={0}
                      maxValue={100}
                      class="w-full"
                    >
                      <Meter.Track class="w-full h-2 bg-bg-muted rounded-full overflow-hidden">
                        <Meter.Fill
                          class={`h-full rounded-full transition-all ease-out duration-150 ${
                            item.status === "error"
                              ? "bg-error"
                              : item.status === "success"
                              ? "bg-success"
                              : "bg-primary"
                          }`}
                          style={{ width: `${item.progress}%` }}
                        />
                      </Meter.Track>
                    </Meter>
                  </div>
                )}
              </For>
            </div>
          </Show>

          <div
            class={`flex ${
              ctx.fileUploads().length > 0 ? "justify-end" : "justify-center"
            } gap-2 mt-6`}
          >
            <Show when={ctx.fileUploads().length === 0}>
              <label class="btn btn-primary cursor-pointer">
                Select Files
                <input
                  type="file"
                  multiple
                  onChange={(e) => ctx.handleFileSelect(e.target.files)}
                  class="hidden"
                />
              </label>
            </Show>

            <Show when={ctx.fileUploads().length > 0}>
              <button
                class={`btn ${
                  ctx.isUploading() ? "btn-disabled" : "btn-outline"
                }`}
                onClick={() => {
                  if (!ctx.isUploading()) {
                    ctx.clearUploads();
                  }
                }}
                disabled={ctx.isUploading()}
              >
                {ctx.isUploading() ? "Uploading..." : "Cancel"}
              </button>
              <button
                class={`btn ${
                  ctx.isUploading() ? "btn-disabled" : "btn-primary"
                }`}
                onClick={() => ctx.handleUploadFiles()}
                disabled={ctx.isUploading() || ctx.fileUploads().length === 0}
              >
                {ctx.isUploading() ? "Uploading..." : "Upload"}
              </button>
            </Show>
          </div>
        </Dialog.Content>
      </Dialog.Portal>
    </Dialog>
  );
}
