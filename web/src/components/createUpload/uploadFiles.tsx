import { Button } from "@kobalte/core/button";
import { Dialog } from "@kobalte/core/dialog";
import { uploadFiles } from "@sv/apis/media";
import { UploadFileInfo } from "@sv/apis/media/models";
import { FileIcon } from "@sv/components/icons";
import useAppCtx from "@sv/store/appCtxProvider";
import FileUtils from "@sv/utils/fileUtils";
import Format from "@sv/utils/format";
import Random from "@sv/utils/random";
import Validate, { VALIDATIONS } from "@sv/utils/validate";
import { createEffect, createSignal, For, Show } from "solid-js";
import { createStore, produce } from "solid-js/store";

// TODO: Get these from the backend
const maxFileSize = 500 * 1024 * 1024; // 500MB
const maxFilesCount = 100;
const maxTotalSize = maxFileSize * maxFilesCount;
const allowedTypes = [
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
  "application/x-rar-compressed",
  "application/x-7z-compressed",
];
// disallowed types has higher priority than allowed types
const disallowedTypes = ["application/pdf"];

interface Props {
  isModalOpen: boolean;
  closeModal: () => void;
}

type UploadFileInfoMap = Record<string, UploadFileInfo>;

export default function UploadFiles(props: Props) {
  let fileInputRef: HTMLInputElement | undefined;

  const appCtx = useAppCtx();

  const [selectedFiles, setSelectedFiles] = createStore<UploadFileInfoMap>({});
  const [isLoading, setIsLoading] = createSignal(false);
  const [error, setError] = createSignal("");
  const [isDragOver, setIsDragOver] = createSignal(false);

  const isDisabled = () =>
    isLoading() ||
    Object.keys(selectedFiles).length === 0 ||
    Object.keys(selectedFiles).length > maxFilesCount;

  const selectedFilesSize = () => {
    return Object.values(selectedFiles).reduce(
      (sum, f) => sum + f.file.size,
      0
    );
  };

  // Reset state when the modal is closed
  createEffect(() => {
    if (!props.isModalOpen) {
      setSelectedFiles({});
      setIsLoading(false);
      setError("");
      setIsDragOver(false);
      fileInputRef = undefined;
    }
  });

  const validateFile = (file: File): string | null => {
    if (!Validate.name(file.name)) {
      return `File name "${file.name}" is invalid. Max length is ${VALIDATIONS.MAX_LENGTH}.`;
    }

    if (file.size > maxFileSize) {
      return `File "${file.name}" is too large. Max size is ${Format.size(
        maxFileSize
      )}.`;
    }

    // Check if file type is allowed
    let isAllowed = false;
    if (allowedTypes.length === 0) {
      isAllowed = true;
    }

    if (!isAllowed) {
      isAllowed = allowedTypes.some((type) => {
        if (type.endsWith("/*")) {
          return file.type.startsWith(type.slice(0, -1));
        }
        return file.type === type;
      });
    }

    if (isAllowed) {
      isAllowed = !disallowedTypes.some((type) => {
        if (type.endsWith("/*")) {
          return file.type.startsWith(type.slice(0, -1));
        }
        return file.type === type;
      });
    }

    if (!isAllowed) {
      return `File type "${file.type}" is not supported.`;
    }

    return null;
  };

  const validateFiles = (files: FileList | File[]): string | null => {
    const newFiles = Array.from(files);

    // Check length
    const existingFiles = Object.values(selectedFiles);
    if (newFiles.length + existingFiles.length > maxFilesCount) {
      return `Too many files selected. Max is ${maxFilesCount} files.`;
    }

    // Process existing files
    const existingFileNames = new Set<string>();
    let totalSize = 0;
    for (const file of existingFiles) {
      existingFileNames.add(file.file.name);
      totalSize += file.file.size;
    }

    const newFileNames = new Set<string>();
    for (const file of newFiles) {
      // Check for duplicates
      if (existingFileNames.has(file.name) || newFileNames.has(file.name)) {
        return `Duplicate file "${file.name}" in the same folder is not allowed.`;
      }
      newFileNames.add(file.name);
      totalSize += file.size;

      // Check total size
      if (totalSize > maxTotalSize) {
        return `Total size of all files is too large. Max is ${Format.size(
          maxTotalSize
        )}.`;
      }
    }

    return null;
  };

  const processFiles = (files: FileList | File[]) => {
    const fileArray = Array.from(files);
    const newFiles: UploadFileInfoMap = {};
    const errs: string[] = [];

    const totalValidationError = validateFiles(files);
    if (totalValidationError) {
      setError(totalValidationError);
      return;
    }

    for (const file of fileArray) {
      const validationError = validateFile(file);
      if (validationError) {
        errs.push(validationError);
        continue;
      }

      const id = Random.id();
      newFiles[id] = {
        id,
        file,
        progress: 0,
        status: "pending",
      };
    }

    if (errs.length > 0) {
      // Show numbered errors
      setError(errs.map((err, i) => `${i + 1}. ${err}`).join("\n"));
      return;
    }

    setError("");
    setSelectedFiles(newFiles);
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
    setSelectedFiles({ [fileId]: undefined });
    setError("");
  };

  const openFileDialog = () => {
    fileInputRef?.click();
  };

  const handleUpload = async () => {
    const files = Object.values(selectedFiles);

    if (files.length === 0) {
      setError("Please select at least one file");
      return;
    }

    setIsLoading(true);
    setError("");

    let errMsg = "";

    const uploadedFiles = uploadFiles(
      files,
      appCtx.currentFolderId(),
      (id, progress) => {
        setSelectedFiles(
          produce((fileMap) => {
            let status: UploadFileInfo["status"] = "pending";
            if (progress === 100) {
              status = "success";
            } else if (progress > 0) {
              status = "uploading";
            }

            fileMap[id].progress = progress;
            fileMap[id].status = status;
          })
        );
      }
    );

    uploadedFiles.forEach(async (fileResult) => {
      try {
        const file = await fileResult.file;
        setSelectedFiles(
          produce((fileMap) => {
            fileMap[fileResult.clientId].status = "success";
          })
        );
      } catch (err) {
        if (!errMsg) {
          errMsg =
            files.length > 1 ? "Failed to upload some files" : "Upload failed";
        }

        setSelectedFiles(
          produce((fileMap) => {
            fileMap[fileResult.clientId].status = "error";
            fileMap[fileResult.clientId].error =
              err instanceof Error ? err.message : "Upload failed";
          })
        );
      }
    });

    setIsLoading(false);

    // If no errors, close modal after brief delay to show success
    if (!errMsg) {
      setTimeout(() => {
        props.closeModal();
      }, 1500);
    }
  };

  const getUploadStats = () => {
    const files = Object.values(selectedFiles);
    const uploading = files.filter((f) => f.status === "uploading").length;
    const success = files.filter((f) => f.status === "success").length;
    const failed = files.filter((f) => f.status === "error").length;

    return { uploading, success, failed, total: files.length };
  };

  return (
    <Dialog
      open={props.isModalOpen}
      onOpenChange={(isOpen) => !isOpen && props.closeModal()}
    >
      <Dialog.Portal>
        <Dialog.Overlay class="dialog-overlay" />
        <Dialog.Content class="dialog-content max-w-2xl">
          <div class="flex flex-col">
            <Dialog.Title class="dialog-title">Upload Files</Dialog.Title>
            <Dialog.Description class="dialog-description">
              Upload your files to the current folder. Maximum 4GB per file, 100
              files total.
            </Dialog.Description>

            <Show when={error()}>
              <div class="input-t-error mb-4">{error()}</div>
            </Show>

            {/* Hidden file input */}
            <input
              ref={fileInputRef}
              type="file"
              multiple
              accept={allowedTypes.join(",")}
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
                  Maximum file size: {Format.size(maxFileSize)} • Maximum files:{" "}
                  {maxFilesCount}
                </p>
              </div>
            </div>

            {/* File list */}
            <Show when={Object.keys(selectedFiles).length > 0}>
              <div class="mt-4">
                <div class="flex justify-between items-center mb-2">
                  <p class="text-sm text-neutral-light">
                    <span class="text-neutral font-medium">
                      {Object.keys(selectedFiles).length}
                    </span>{" "}
                    file(s) selected • Total size:{" "}
                    {Format.size(selectedFilesSize())}
                  </p>

                  <Show when={isLoading()}>
                    <div class="text-xs text-neutral-light">
                      {(() => {
                        const stats = getUploadStats();
                        return `${stats.success}/${stats.total} completed`;
                      })()}
                    </div>
                  </Show>
                </div>

                <div class="max-h-[300px] md:max-h-[400px] overflow-y-auto space-y-2 mb-4">
                  <For each={Object.values(selectedFiles)}>
                    {(fileWithId) => (
                      <div class="flex items-center justify-between p-3 border border-border rounded bg-white">
                        <div class="flex items-center gap-3 flex-1 min-w-0">
                          {/* File preview/icon */}
                          <div class="size-10 md:size-12 rounded bg-bg-muted flex items-center justify-center flex-shrink-0">
                            <Show
                              when={
                                fileWithId.file.type.startsWith("image/") &&
                                fileWithId.file.size < 10 * 1024 * 1024
                              }
                              fallback={
                                <FileIcon
                                  fileCategory={FileUtils.mimeToCategory(
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
                              {Format.size(fileWithId.file.size)}
                              <Show
                                when={fileWithId.file.size > 100 * 1024 * 1024}
                              >
                                <span class="ml-1 text-blue-600">
                                  • Chunked upload
                                </span>
                              </Show>
                            </div>

                            {/* Progress bar for uploading files */}
                            <Show when={fileWithId.status === "uploading"}>
                              <div class="w-full bg-bg-muted rounded-full h-1.5 mt-1">
                                <div
                                  class="bg-primary h-1.5 rounded-full transition-all duration-300"
                                  style={{ width: `${fileWithId.progress}%` }}
                                />
                              </div>
                              <div class="text-xs text-neutral-light mt-1">
                                {Math.round(fileWithId.progress)}%
                              </div>
                            </Show>

                            {/* Status indicators */}
                            <Show when={fileWithId.status === "success"}>
                              <div class="text-xs text-green-600 mt-1">
                                ✓ Uploaded successfully
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
                {isLoading()
                  ? "Uploading..."
                  : `Upload ${Object.keys(selectedFiles).length} File(s)`}
              </Button>
            </div>
          </div>
        </Dialog.Content>
      </Dialog.Portal>
    </Dialog>
  );
}
