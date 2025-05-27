import { Dialog } from "@kobalte/core/dialog";
import { Meter } from "@kobalte/core/meter";
import { FileField } from "@kobalte/core/file-field";
import { createEffect, createSignal, For, Show } from "solid-js";
import { Button } from "@kobalte/core/button";
import { FileIcon } from "@sv/utils/icons";
import { FILE_CATEGORIES } from "@sv/utils/consts";

// Handle the FileField details type more generically
// type Details = { value: File[] };

interface Props {
  isModalOpen: boolean;
  closeModal: () => void;
}

export default function UploadFiles(props: Props) {
  let formRef!: HTMLFormElement;

  // TODO: Instead of form and formRef, manually track the files using the handleFileChange which receives this structure: {acceptedFiles: Array<File>, rejectedFiles: Array<File>}
  const [files, setFiles] = createSignal<File[]>([]);
  const [isLoading, setIsLoading] = createSignal(false);
  const [error, setError] = createSignal("");

  const isDisabled = () => isLoading();

  // Reset files when modal closes
  // createEffect(() => {
  //   if (!ctx.isFileUploadModalOpen()) {
  //     setFiles([]);
  //   }
  // });

  const handleFileChange = (data: any) => {
    console.log("handleFileChange: ", data);
  };

  const handleReset = () => {
    setError("");
    new FormData(formRef).delete("selected-files");
    props.closeModal();
  };

  const handleStartUpload = () => {
    const formData = new FormData(formRef);
    const selectedFiles = formData.getAll("selected-files");

    if (selectedFiles.length === 0) {
      setError("Please select at least one file");
      return;
    }

    console.log(selectedFiles);

    // Convert files array to FileList-like object for our context
    // const fileList = Object.assign(files(), {
    //   item: (i: number) => files()[i],
    // }) as unknown as FileList;

    // // Pass to context for upload
    // ctx.handleFileSelect(fileList);
    // ctx.handleUploadFiles();
  };

  // This will show a different file icon based on mime type
  // const getFileTypeIcon = (type: string) => {
  //   if (type.startsWith("image/")) {
  //     return (
  //       <svg
  //         xmlns="http://www.w3.org/2000/svg"
  //         fill="none"
  //         viewBox="0 0 24 24"
  //         stroke-width="1.5"
  //         stroke="currentColor"
  //         class="size-6"
  //       >
  //         <path
  //           stroke-linecap="round"
  //           stroke-linejoin="round"
  //           d="m2.25 15.75 5.159-5.159a2.25 2.25 0 0 1 3.182 0l5.159 5.159m-1.5-1.5 1.409-1.409a2.25 2.25 0 0 1 3.182 0l2.909 2.909m-18 3.75h16.5a1.5 1.5 0 0 0 1.5-1.5V6a1.5 1.5 0 0 0-1.5-1.5H3.75A1.5 1.5 0 0 0 2.25 6v12a1.5 1.5 0 0 0 1.5 1.5Zm10.5-11.25h.008v.008h-.008V8.25Zm.375 0a.375.375 0 1 1-.75 0 .375.375 0 0 1 .75 0Z"
  //         />
  //       </svg>
  //     );
  //   } else if (type.startsWith("video/")) {
  //     return (
  //       <svg
  //         xmlns="http://www.w3.org/2000/svg"
  //         fill="none"
  //         viewBox="0 0 24 24"
  //         stroke-width="1.5"
  //         stroke="currentColor"
  //         class="size-6"
  //       >
  //         <path
  //           stroke-linecap="round"
  //           stroke-linejoin="round"
  //           d="M3.375 19.5h17.25m-17.25 0a1.125 1.125 0 0 1-1.125-1.125M3.375 19.5h1.5C5.496 19.5 6 18.996 6 18.375m-3.75 0V5.625m0 12.75v-1.5c0-.621.504-1.125 1.125-1.125m18.375 2.625V5.625m0 12.75c0 .621-.504 1.125-1.125 1.125m1.125-1.125v-1.5c0-.621-.504-1.125-1.125-1.125m0 3.75h-1.5A1.125 1.125 0 0 1 18 18.375M20.625 4.5H3.375m17.25 0c.621 0 1.125.504 1.125 1.125M20.625 4.5h-1.5C18.504 4.5 18 5.004 18 5.625m3.75 0v1.5c0 .621-.504 1.125-1.125 1.125M3.375 4.5c-.621 0-1.125.504-1.125 1.125M3.375 4.5h1.5C5.496 4.5 6 5.004 6 5.625m-3.75 0v1.5c0 .621.504 1.125 1.125 1.125m0 0h1.5m-1.5 0c-.621 0-1.125.504-1.125 1.125v1.5c0 .621.504 1.125 1.125 1.125m1.5-3.75C5.496 8.25 6 7.746 6 7.125v-1.5M4.875 8.25C5.496 8.25 6 8.754 6 9.375v1.5m0-5.25v5.25m0-5.25C6 5.004 6.504 4.5 7.125 4.5h9.75c.621 0 1.125.504 1.125 1.125m1.125 2.625h1.5m-1.5 0A1.125 1.125 0 0 1 18 7.125v-1.5m1.125 2.625c-.621 0-1.125.504-1.125 1.125v1.5m2.625-2.625c.621 0 1.125.504 1.125 1.125v1.5c0 .621-.504 1.125-1.125 1.125M18 5.625v5.25M7.125 12h9.75m-9.75 0A1.125 1.125 0 0 1 6 10.875M7.125 12C6.504 12 6 12.504 6 13.125m0-2.25C6 11.496 5.496 12 4.875 12M18 10.875c0 .621-.504 1.125-1.125 1.125M18 10.875c0 .621.504 1.125 1.125 1.125m-2.25 0c.621 0 1.125.504 1.125 1.125m-12 5.25v-5.25m0 5.25c0 .621.504 1.125 1.125 1.125h9.75c.621 0 1.125-.504 1.125-1.125m-12 0v-1.5c0-.621-.504-1.125-1.125-1.125M18 18.375v-5.25m0 5.25v-1.5c0-.621.504-1.125 1.125-1.125M18 13.125v1.5c0 .621.504 1.125 1.125 1.125M18 13.125c0 .621.504 1.125 1.125 1.125m-1.5-1.5v-5.25m0 5.25v-1.5c0-.621-.504-1.125-1.125-1.125m0 0h-9.75M18 13.125c0 .621-.504 1.125-1.125 1.125M6 13.125v1.5c0 .621-.504 1.125-1.125 1.125M6 13.125C6 12.504 5.496 12 4.875 12m-1.5 0h1.5m-1.5 0c-.621 0-1.125.504-1.125 1.125v1.5c0 .621.504 1.125 1.125 1.125M19.125 12h1.5m0 0c.621 0 1.125.504 1.125 1.125v1.5c0 .621-.504 1.125-1.125 1.125m-17.25 0h1.5m14.25 0h1.5"
  //         />
  //       </svg>
  //     );
  //   } else if (type.startsWith("audio/")) {
  //     return (
  //       <svg
  //         xmlns="http://www.w3.org/2000/svg"
  //         fill="none"
  //         viewBox="0 0 24 24"
  //         stroke-width="1.5"
  //         stroke="currentColor"
  //         class="size-6"
  //       >
  //         <path
  //           stroke-linecap="round"
  //           stroke-linejoin="round"
  //           d="M9 9l10.5-3m0 6.553v3.75a2.25 2.25 0 0 1-1.632 2.163l-1.32.377a1.803 1.803 0 1 1-.99-3.467l2.31-.66a2.25 2.25 0 0 0 1.632-2.163Zm0 0V2.25L9 5.25v10.303m0 0v3.75a2.25 2.25 0 0 1-1.632 2.163l-1.32.377a1.803 1.803 0 0 1-.99-3.467l2.31-.66A2.25 2.25 0 0 0 9 15.553Z"
  //         />
  //       </svg>
  //     );
  //   } else if (type === "application/pdf") {
  //     return (
  //       <svg
  //         xmlns="http://www.w3.org/2000/svg"
  //         fill="none"
  //         viewBox="0 0 24 24"
  //         stroke-width="1.5"
  //         stroke="currentColor"
  //         class="size-6"
  //       >
  //         <path
  //           stroke-linecap="round"
  //           stroke-linejoin="round"
  //           d="M19.5 14.25v-2.625a3.375 3.375 0 0 0-3.375-3.375h-1.5A1.125 1.125 0 0 1 13.5 7.125v-1.5a3.375 3.375 0 0 0-3.375-3.375H8.25m2.25 0H5.625c-.621 0-1.125.504-1.125 1.125v17.25c0 .621.504 1.125 1.125 1.125h12.75c.621 0 1.125-.504 1.125-1.125V11.25a9 9 0 0 0-9-9Z"
  //         />
  //       </svg>
  //     );
  //   } else {
  //     return (
  //       <svg
  //         xmlns="http://www.w3.org/2000/svg"
  //         fill="none"
  //         viewBox="0 0 24 24"
  //         stroke-width="1.5"
  //         stroke="currentColor"
  //         class="size-6"
  //       >
  //         <path
  //           stroke-linecap="round"
  //           stroke-linejoin="round"
  //           d="M19.5 14.25v-2.625a3.375 3.375 0 0 0-3.375-3.375h-1.5A1.125 1.125 0 0 1 13.5 7.125v-1.5a3.375 3.375 0 0 0-3.375-3.375H8.25m2.25 0H5.625c-.621 0-1.125.504-1.125 1.125v17.25c0 .621.504 1.125 1.125 1.125h12.75c.621 0 1.125-.504 1.125-1.125V11.25a9 9 0 0 0-9-9Z"
  //         />
  //       </svg>
  //     );
  //   }
  // };

  return (
    <Dialog
      open={props.isModalOpen}
      onOpenChange={(isOpen) => !isOpen && handleReset()}
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

            {/* <form
              ref={formRef}
              onSubmit={(e) => {
                e.preventDefault();
                handleStartUpload();
              }}
              onReset={(e) => {
                e.preventDefault();
                handleReset();
              }}
            >
              <FileField
                multiple
                allowDragAndDrop
                maxFileSize={1 * 1024 * 1024} // 1MB max per file
                // onFileChange={handleFileChange}
                // onFileAccept={handleFileAccept}
                // onFileReject={handleFileReject}
                disabled={isLoading()}
              >
                <FileField.Dropzone>
                  Drop your files here...
                  <FileField.Trigger>Choose files</FileField.Trigger>
                </FileField.Dropzone>

                <FileField.HiddenInput name="selected-files" />

                <FileField.ErrorMessage class="text-sm text-error mt-2" />

                <FileField.ItemList>
                  {(f) => (
                    <FileField.Item>
                      <FileField.ItemPreviewImage />
                      <FileField.ItemName />
                      <FileField.ItemSize />
                      <FileField.ItemDeleteTrigger>
                        Delete
                      </FileField.ItemDeleteTrigger>
                    </FileField.Item>
                  )}
                </FileField.ItemList>
              </FileField>

              <div class="flex justify-end gap-2 mt-4">
                <Button
                  type="reset"
                  class="btn btn-outline"
                  onClick={props.closeModal}
                  disabled={isLoading()}
                >
                  Cancel
                </Button>
                <Button
                  type="submit"
                  classList={{
                    btn: true,
                    "btn-disabled": isDisabled(),
                    "btn-primary": !isDisabled(),
                  }}
                  disabled={isDisabled()}
                >
                  {isLoading() ? "Uploading..." : "Upload Files"}
                </Button>
              </div>
            </form> */}

            <form
              ref={formRef}
              onSubmit={(e) => {
                e.preventDefault();
                handleStartUpload();
              }}
              onReset={(e) => {
                e.preventDefault();
                handleReset();
              }}
            >
              <FileField
                multiple
                maxFiles={2}
                allowDragAndDrop
                maxFileSize={1 * 1024 * 1024} // 1MB max per file
                onFileChange={handleFileChange}
                disabled={isLoading()}
              >
                {/* <FileField.Label class="block mb-2 text-neutral-light">
                  Upload your files
                </FileField.Label> */}

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
                      <span class="text-neutral font-medium">
                        click anywhere
                      </span>{" "}
                      to select files
                    </p>
                  </div>
                </FileField.Dropzone>

                <FileField.HiddenInput name="selected-files" />

                <div class="mt-4">
                  <FileField.Description class="text-sm text-neutral-light mb-2">
                    {new FormData(formRef).getAll("selected-files").length}
                    file(s) selected. Click "Upload Files" to start the upload
                    process.
                  </FileField.Description>

                  <FileField.ItemList class="max-h-[200px] overflow-y-auto space-y-2 mb-4">
                    {(file) => (
                      <FileField.Item class="flex items-center justify-between p-2 border border-border rounded">
                        <div class="flex items-center gap-2">
                          <FileField.ItemPreview
                            type={file.type}
                            class="w-8 h-8 rounded bg-bg-muted flex items-center justify-center"
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

                          <div>
                            <FileField.ItemName class="font-medium text-sm truncate max-w-[180px]" />
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
                <FileField.ErrorMessage class="text-sm text-error mt-2" />
              </FileField>
              {/* <div class="max-h-[300px] overflow-y-auto mb-4">
                <For each={ctx.fileUploads()}>
                  {(item) => (
                    <div class="mb-4 pb-4 border-b border-border last:border-b-0 last:pb-0">
                      <div class="flex justify-between mb-1">
                        <div class="flex items-center gap-2">
                          <div class="rounded bg-bg-muted w-8 h-8 flex items-center justify-center text-neutral-light overflow-hidden">
                            {item.file.type.startsWith("image/") ? (
                              <img
                                src={URL.createObjectURL(item.file)}
                                alt={item.file.name}
                                class="w-full h-full object-cover"
                                onLoad={(e) =>
                                  URL.revokeObjectURL(e.currentTarget.src)
                                }
                              />
                            ) : (
                              getFileTypeIcon(item.file.type)
                            )}
                          </div>
                          <div>
                            <div
                              class="truncate max-w-[200px] font-medium text-sm"
                              title={item.file.name}
                            >
                              {item.file.name}
                            </div>
                            <div class="text-xs text-neutral-light">
                              {(item.file.size / 1024).toFixed(1)} KB
                            </div>
                          </div>
                        </div>
                        <div
                          class={`text-sm ${
                            item.status === "success"
                              ? "text-success"
                              : item.status === "error"
                              ? "text-error"
                              : "text-neutral"
                          }`}
                        >
                          {item.status === "pending" && "Pending"}
                          {item.status === "uploading" && `${item.progress}%`}
                          {item.status === "success" && "Completed"}
                          {item.status === "error" && "Failed"}
                        </div>
                      </div>

                      <div class="ml-10 mt-2">
                        <Meter
                          value={item.progress}
                          minValue={0}
                          maxValue={100}
                          class="w-full"
                          getValueLabel={({ value }) => `${value}% complete`}
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
                    </div>
                  )}
                </For>

                <div class="flex justify-end gap-2 mt-6">
                  <button
                    class="btn btn-outline"
                    onClick={() => ctx.setIsFileUploadModalOpen(false)}
                    disabled={ctx.isUploading()}
                  >
                    Close
                  </button>
                </div>
              </div> */}
              <div class="flex justify-end gap-2 mt-4">
                <Button
                  type="reset"
                  class="btn btn-outline"
                  onClick={props.closeModal}
                  disabled={isLoading()}
                >
                  Cancel
                </Button>
                <Button
                  type="submit"
                  classList={{
                    btn: true,
                    "btn-disabled": isDisabled(),
                    "btn-primary": !isDisabled(),
                  }}
                  disabled={isDisabled()}
                >
                  {isLoading() ? "Uploading..." : "Upload Files"}
                </Button>
              </div>
            </form>
          </div>
        </Dialog.Content>
      </Dialog.Portal>
    </Dialog>
  );
}
