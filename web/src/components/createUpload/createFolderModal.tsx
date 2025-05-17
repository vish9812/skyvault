import { Dialog } from "@kobalte/core/dialog";
import { createEffect } from "solid-js";

interface CreateFolderModalProps {
  isOpen: boolean;
  onClose: () => void;
  folderName: string;
  setFolderName: (name: string) => void;
  onCreateFolder: () => void;
  isCreating: boolean;
  error: string | null;
}

export function CreateFolderModal(props: CreateFolderModalProps) {
  // Reset the folder name when the modal opens
  createEffect(() => {
    if (props.isOpen) {
      props.setFolderName("");
    }
  });

  return (
    <Dialog
      open={props.isOpen}
      onOpenChange={(isOpen) => !isOpen && props.onClose()}
    >
      <Dialog.Portal>
        <Dialog.Overlay class="fixed inset-0 bg-black/20 z-50" />
        <Dialog.Content class="fixed left-1/2 top-1/2 -translate-x-1/2 -translate-y-1/2 bg-white rounded-lg p-6 shadow-xl z-50 w-full max-w-md">
          <div class="flex flex-col gap-4">
            <Dialog.Title class="text-xl font-bold">
              Create New Folder
            </Dialog.Title>
            <Dialog.Description class="subtitle text-neutral-light">
              Enter a name for your new folder
            </Dialog.Description>

            <div class="mt-2">
              <label for="folder-name" class="label">
                Folder Name
              </label>
              <input
                id="folder-name"
                class="input input-b-std w-full"
                type="text"
                value={props.folderName}
                onInput={(e) => props.setFolderName(e.currentTarget.value)}
                placeholder="Enter folder name"
                autofocus
              />
              {props.error && <p class="input-t-error mt-1">{props.error}</p>}
            </div>

            <div class="flex justify-end gap-2 mt-4">
              <button
                class="btn btn-outline"
                onClick={props.onClose}
                disabled={props.isCreating}
              >
                Cancel
              </button>
              <button
                class="btn btn-primary"
                onClick={props.onCreateFolder}
                disabled={props.isCreating}
              >
                {props.isCreating ? "Creating..." : "Create Folder"}
              </button>
            </div>
          </div>
        </Dialog.Content>
      </Dialog.Portal>
    </Dialog>
  );
}
