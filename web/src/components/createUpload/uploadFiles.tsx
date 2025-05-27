import { Dialog } from "@kobalte/core/dialog";
import { Meter } from "@kobalte/core/meter";
import { Details, FileField } from "@kobalte/core/file-field";
import { createEffect, createSignal, For, Show } from "solid-js";
import { Button } from "@kobalte/core/button";
import { FileIcon } from "@sv/utils/icons";
import { FILE_CATEGORIES } from "@sv/utils/consts";

interface Props {
  isModalOpen: boolean;
  closeModal: () => void;
}

export default function UploadFiles(props: Props) {
  const [acceptedFiles, setAcceptedFiles] = createSignal<File[]>([]);
  const [isLoading, setIsLoading] = createSignal(false);
  const [error, setError] = createSignal("");

  const isDisabled = () => isLoading() || acceptedFiles().length === 0;

  // Reset files when modal closes
  createEffect(() => {
    if (!props.isModalOpen) {
      setAcceptedFiles([]);
      setError("");
      setIsLoading(false);
    }
  });

  const handleFileChange = (details: Details) => {
    console.log("File change details: ", details);
    console.log("Current accepted files: ", acceptedFiles());

    // const currentAccepted = acceptedFiles();
    // const newAccepted = details.acceptedFiles;
    // const newRejected = details.rejectedFiles;

    // // Case 1: New files selected (new files are added to existing ones)
    // // This happens when user selects additional files
    // if (newAccepted.length > currentAccepted.length) {
    //   // Add only the new files that aren't already in the accepted list
    //   const filesToAdd = newAccepted.filter(
    //     (newFile) =>
    //       !currentAccepted.some(
    //         (existingFile) =>
    //           existingFile.name === newFile.name &&
    //           existingFile.size === newFile.size &&
    //           existingFile.lastModified === newFile.lastModified
    //       )
    //   );
    //   setAcceptedFiles([...currentAccepted, ...filesToAdd]);
    // }
    // // Case 2: File removed (details.acceptedFiles contains all remaining files)
    // // This happens when user clicks the delete trigger
    // else if (newAccepted.length < currentAccepted.length) {
    //   setAcceptedFiles(newAccepted);
    // }
    // // Case 3: Files rejected (details.acceptedFiles might be empty, focus on rejectedFiles)
    // // This happens when files don't meet the criteria (size, type, etc.)
    // // else if (newRejected.length > rejectedFiles().length) {
    // //   // Keep existing accepted files, just update rejected files
    // //   setAcceptedFiles(currentAccepted);
    // // }
    // // Edge case: Same number of files but different files (replacement)
    // else {
    //   setAcceptedFiles(newAccepted);
    // }

    // // Always update rejected files
    // // setRejectedFiles(newRejected);

    // // Clear error when files change
    // setError("");

    console.log("Final accepted files: ", acceptedFiles());
  };

  const handleUpload = () => {
    const files = acceptedFiles();

    if (files.length === 0) {
      setError("Please select at least one file");
      return;
    }

    console.log("Starting upload with files:", files);

    // TODO: Implement actual upload logic here
    // You can now use the files array directly for upload
    // Example:
    // setIsLoading(true);
    // try {
    //   await uploadFiles(files);
    //   props.closeModal();
    // } catch (error) {
    //   setError("Upload failed");
    // } finally {
    //   setIsLoading(false);
    // }
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
              <div class="input-t-error">{error()}</div>
            </Show>

            <FileField
              multiple
              allowDragAndDrop
              maxFileSize={1 * 1024 * 1024} // 1MB max per file
              onFileChange={handleFileChange}
              disabled={isLoading()}
            >
              <FileField.Dropzone class="border-2 border-dashed rounded-lg p-6 my-4 text-center transition-colors hover:border-primary hover:bg-primary-lighter hover:cursor-pointer">
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
                    Or{" "}
                    <span class="text-neutral font-medium">click anywhere</span>{" "}
                    to select files
                  </p>
                </div>
              </FileField.Dropzone>
              <FileField.HiddenInput />
              <div class="mt-4">
                <FileField.Description class="text-sm text-neutral-light mb-2">
                  <span class="text-neutral font-medium">
                    {acceptedFiles().length}
                  </span>{" "}
                  file(s) selected. Click "Upload Files" to start the upload
                  process.
                </FileField.Description>

                <FileField.ItemList class="max-h-[200px] md:max-h-[300px] overflow-y-auto space-y-2 mb-4">
                  {(file) => (
                    <FileField.Item class="flex items-center justify-between p-2 border border-border rounded">
                      <div class="flex items-center gap-2">
                        <FileField.ItemPreview
                          type={file.type}
                          class="size-8 md:size-12 rounded bg-bg-muted flex items-center justify-center"
                        >
                          <Show
                            when={file.type.startsWith("image/")}
                            fallback={
                              <FileIcon
                                fileCategory={
                                  file.type.split("/")[0] as FILE_CATEGORIES
                                }
                                isFolder={false}
                              />
                            }
                          >
                            <FileField.ItemPreviewImage class="w-full h-full object-cover rounded" />
                          </Show>
                        </FileField.ItemPreview>

                        <div class="flex flex-col">
                          <FileField.ItemName class="font-medium text-sm truncate max-w-[200px] md:max-w-[300px]" />
                          <FileField.ItemSize
                            class="text-xs text-neutral-light"
                            precision={1}
                          />
                        </div>
                      </div>

                      <FileField.ItemDeleteTrigger class="text-neutral-light hover:text-error p-1 rounded">
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
                      </FileField.ItemDeleteTrigger>
                    </FileField.Item>
                  )}
                </FileField.ItemList>
              </div>
            </FileField>

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
