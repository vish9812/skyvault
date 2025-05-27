import { Dialog } from "@kobalte/core/dialog";
import { createEffect, createSignal, For, Show } from "solid-js";
import { Button } from "@kobalte/core/button";
import { FileIcon } from "@sv/utils/icons";
import { FILE_CATEGORIES, ROOT_FOLDER_ID } from "@sv/utils/consts";
import { uploadFiles } from "@sv/apis/media";
import useAppCtx from "@sv/store/appCtxProvider";
import format from "@sv/utils/format";

interface Props {
  isModalOpen: boolean;
  closeModal: () => void;
}

interface FileWithId {
  id: string;
  file: File;
  progress: number;
  status: "pending" | "uploading" | "success" | "error";
  error?: string;
}

export default function UploadFiles(props: Props) {
  const appCtx = useAppCtx();
  const [selectedFiles, setSelectedFiles] = createSignal<FileWithId[]>([]);
  const [isLoading, setIsLoading] = createSignal(false);
  const [error, setError] = createSignal("");
  const [isDragOver, setIsDragOver] = createSignal(false);

  let fileInputRef: HTMLInputElement | undefined;

  const MAX_FILE_SIZE = 10 * 1024 * 1024; // 10MB
  const ALLOWED_TYPES = [
    "image/*",
    "video/*",
    "audio/*",
    "text/*",
    "application/pdf",
    "application/msword",
    "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
    "application/vnd.ms-excel",
    "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
    "application/zip",
    "application/x-zip-compressed",
  ];

  const isDisabled = () => isLoading() || selectedFiles().length === 0;

  // Reset files when modal closes
  createEffect(() => {
    if (!props.isModalOpen) {
      setSelectedFiles([]);
      setError("");
      setIsLoading(false);
      setIsDragOver(false);
      if (fileInputRef) {
        fileInputRef.value = "";
      }
    }
  });

  const generateFileId = () => Math.random().toString(36).substr(2, 9);

  const getFileCategory = (mimeType: string): FILE_CATEGORIES => {
    const type = mimeType.split("/")[0];
    switch (type) {
      case "image":
        return FILE_CATEGORIES.IMAGE;
      case "video":
        return FILE_CATEGORIES.VIDEO;
      case "audio":
        return FILE_CATEGORIES.AUDIO;
      case "text":
        return FILE_CATEGORIES.TEXT;
      default:
        return FILE_CATEGORIES.OTHER;
    }
  };

  const validateFile = (file: File): string | null => {
    if (file.size > MAX_FILE_SIZE) {
      return `File "${file.name}" is too large. Maximum size is ${format.size(
        MAX_FILE_SIZE
      )}.`;
    }

    // Check if file type is allowed
    const isAllowed = ALLOWED_TYPES.some((type) => {
      if (type.endsWith("/*")) {
        return file.type.startsWith(type.slice(0, -1));
      }
      return file.type === type;
    });

    if (!isAllowed) {
      return `File type "${file.type}" is not supported.`;
    }

    return null;
  };

  const processFiles = (files: FileList | File[]) => {
    const fileArray = Array.from(files);
    const newFiles: FileWithId[] = [];
    const errors: string[] = [];

    fileArray.forEach((file) => {
      const validationError = validateFile(file);
      if (validationError) {
        errors.push(validationError);
        return;
      }

      // Check for duplicates
      const isDuplicate = selectedFiles().some(
        (existing) =>
          existing.file.name === file.name &&
          existing.file.size === file.size &&
          existing.file.lastModified === file.lastModified
      );

      if (!isDuplicate) {
        newFiles.push({
          id: generateFileId(),
          file,
          progress: 0,
          status: "pending",
        });
      }
    });

    if (errors.length > 0) {
      setError(errors.join(" "));
    } else {
      setError("");
    }

    if (newFiles.length > 0) {
      setSelectedFiles((prev) => [...prev, ...newFiles]);
    }
  };

  const handleFileInputChange = (event: Event) => {
    const target = event.target as HTMLInputElement;
    if (target.files && target.files.length > 0) {
      processFiles(target.files);
    }
  };

  const handleDragOver = (event: DragEvent) => {
    event.preventDefault();
    event.stopPropagation();
    setIsDragOver(true);
  };

  const handleDragLeave = (event: DragEvent) => {
    event.preventDefault();
    event.stopPropagation();
    setIsDragOver(false);
  };

  const handleDrop = (event: DragEvent) => {
    event.preventDefault();
    event.stopPropagation();
    setIsDragOver(false);

    const files = event.dataTransfer?.files;
    if (files && files.length > 0) {
      processFiles(files);
    }
  };

  const removeFile = (fileId: string) => {
    setSelectedFiles((prev) => prev.filter((f) => f.id !== fileId));
    setError("");
  };

  const openFileDialog = () => {
    fileInputRef?.click();
  };

  const createImagePreview = (file: File): Promise<string> => {
    return new Promise((resolve, reject) => {
      if (!file.type.startsWith("image/")) {
        reject(new Error("Not an image file"));
        return;
      }

      const reader = new FileReader();
      reader.onload = (e) => {
        resolve(e.target?.result as string);
      };
      reader.onerror = reject;
      reader.readAsDataURL(file);
    });
  };

  const handleUpload = async () => {
    const files = selectedFiles();

    if (files.length === 0) {
      setError("Please select at least one file");
      return;
    }

    setIsLoading(true);
    setError("");

    try {
      // Mark all files as uploading
      setSelectedFiles((prev) =>
        prev.map((item) => ({ ...item, status: "uploading" as const }))
      );

      // Convert to simple File array for the API
      const fileArray = files.map((f) => f.file);

      await uploadFiles(
        fileArray,
        appCtx.currentFolderId() || ROOT_FOLDER_ID,
        (fileIndex, progress) => {
          setSelectedFiles((prev) => {
            const updated = [...prev];
            if (updated[fileIndex]) {
              updated[fileIndex] = {
                ...updated[fileIndex],
                progress,
              };
            }
            return updated;
          });
        }
      );

      // Mark all as success
      setSelectedFiles((prev) =>
        prev.map((item) => ({
          ...item,
          status: "success" as const,
          progress: 100,
        }))
      );

      // Close modal after brief delay to show success
      setTimeout(() => {
        props.closeModal();
      }, 1500);
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : "Upload failed";
      setError(errorMessage);

      // Mark files as failed
      setSelectedFiles((prev) =>
        prev.map((item) =>
          item.status === "uploading"
            ? { ...item, status: "error" as const, error: errorMessage }
            : item
        )
      );
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <Dialog
      open={props.isModalOpen}
      onOpenChange={(isOpen) => !isOpen && props.closeModal()}
    >
      <Dialog.Portal>
        <Dialog.Overlay class="dialog-overlay" />
        <Dialog.Content class="dialog-content max-w-lg">
          <div class="flex flex-col">
            <Dialog.Title class="dialog-title">Upload Files</Dialog.Title>
            <Dialog.Description class="dialog-description">
              Upload your files to the current folder
            </Dialog.Description>

            <Show when={error()}>
              <div class="input-t-error mb-4">{error()}</div>
            </Show>

            {/* Hidden file input */}
            <input
              ref={fileInputRef}
              type="file"
              multiple
              accept={ALLOWED_TYPES.join(",")}
              onChange={handleFileInputChange}
              style={{ display: "none" }}
              disabled={isLoading()}
            />

            {/* Drag and drop zone */}
            <div
              class={`border-2 border-dashed rounded-lg p-6 my-4 text-center transition-all cursor-pointer ${
                isDragOver()
                  ? "border-primary bg-primary-lighter"
                  : "border-border hover:border-primary hover:bg-primary-lighter"
              } ${isLoading() ? "opacity-50 cursor-not-allowed" : ""}`}
              onDragOver={handleDragOver}
              onDragLeave={handleDragLeave}
              onDrop={handleDrop}
              onClick={!isLoading() ? openFileDialog : undefined}
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
                <p class="font-medium">
                  {isDragOver()
                    ? "Drop files here"
                    : "Drag and drop files here"}
                </p>
                <p class="text-sm text-neutral-light">
                  Or{" "}
                  <span class="text-neutral font-medium">click anywhere</span>{" "}
                  to select files
                </p>
                <p class="text-xs text-neutral-light mt-1">
                  Maximum file size: {format.size(MAX_FILE_SIZE)}
                </p>
              </div>
            </div>

            {/* File list */}
            <Show when={selectedFiles().length > 0}>
              <div class="mt-4">
                <p class="text-sm text-neutral-light mb-2">
                  <span class="text-neutral font-medium">
                    {selectedFiles().length}
                  </span>{" "}
                  file(s) selected
                </p>

                <div class="max-h-[200px] md:max-h-[300px] overflow-y-auto space-y-2 mb-4">
                  <For each={selectedFiles()}>
                    {(fileWithId) => (
                      <div class="flex items-center justify-between p-3 border border-border rounded bg-white">
                        <div class="flex items-center gap-3 flex-1 min-w-0">
                          {/* File preview/icon */}
                          <div class="size-10 md:size-12 rounded bg-bg-muted flex items-center justify-center flex-shrink-0">
                            <Show
                              when={fileWithId.file.type.startsWith("image/")}
                              fallback={
                                <FileIcon
                                  fileCategory={getFileCategory(
                                    fileWithId.file.type
                                  )}
                                  isFolder={false}
                                  size={6}
                                />
                              }
                            >
                              <img
                                src={URL.createObjectURL(fileWithId.file)}
                                alt={fileWithId.file.name}
                                class="w-full h-full object-cover rounded"
                                onLoad={(e) => {
                                  // Clean up object URL after image loads
                                  setTimeout(() => {
                                    URL.revokeObjectURL(
                                      (e.target as HTMLImageElement).src
                                    );
                                  }, 1000);
                                }}
                              />
                            </Show>
                          </div>

                          {/* File info */}
                          <div class="flex flex-col min-w-0 flex-1">
                            <div class="font-medium text-sm truncate">
                              {fileWithId.file.name}
                            </div>
                            <div class="text-xs text-neutral-light">
                              {format.size(fileWithId.file.size)}
                            </div>

                            {/* Progress bar for uploading files */}
                            <Show when={fileWithId.status === "uploading"}>
                              <div class="w-full bg-bg-muted rounded-full h-1.5 mt-1">
                                <div
                                  class="bg-primary h-1.5 rounded-full transition-all duration-300"
                                  style={{ width: `${fileWithId.progress}%` }}
                                />
                              </div>
                            </Show>

                            {/* Status indicators */}
                            <Show when={fileWithId.status === "success"}>
                              <div class="text-xs text-green-600 mt-1">
                                ✓ Uploaded
                              </div>
                            </Show>
                            <Show when={fileWithId.status === "error"}>
                              <div class="text-xs text-red-600 mt-1">
                                ✗ {fileWithId.error || "Upload failed"}
                              </div>
                            </Show>
                          </div>
                        </div>

                        {/* Remove button */}
                        <Show when={fileWithId.status !== "uploading"}>
                          <button
                            class="text-neutral-light hover:text-error p-1 rounded ml-2 flex-shrink-0"
                            onClick={() => removeFile(fileWithId.id)}
                            disabled={isLoading()}
                          >
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
                          </button>
                        </Show>
                      </div>
                    )}
                  </For>
                </div>
              </div>
            </Show>

            {/* Action buttons */}
            <div class="flex justify-end gap-2 mt-4">
              <Button
                class="btn btn-outline"
                onClick={props.closeModal}
                disabled={isLoading()}
              >
                Cancel
              </Button>
              <Button
                classList={{
                  btn: true,
                  "btn-disabled": isDisabled(),
                  "btn-primary": !isDisabled(),
                }}
                disabled={isDisabled()}
                onClick={handleUpload}
              >
                {isLoading() ? "Uploading..." : "Upload Files"}
              </Button>
            </div>
          </div>
        </Dialog.Content>
      </Dialog.Portal>
    </Dialog>
  );
}
