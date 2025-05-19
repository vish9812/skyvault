import { Dialog } from "@kobalte/core/dialog";
import useCtx from "./ctxProvider";
import { useParams } from "@solidjs/router";

export function CreateFolder() {
  const ctx = useCtx();
  const params = useParams();

  return (
    <Dialog
      open={ctx.isFolderModalOpen()}
      onOpenChange={(isOpen) => !isOpen && ctx.setIsFolderModalOpen(false)}
    >
      <Dialog.Portal>
        <Dialog.Overlay class="dialog-overlay" />
        <Dialog.Content class="dialog-content">
          <div class="flex flex-col">
            <Dialog.Title class="dialog-title">Create New Folder</Dialog.Title>
            <Dialog.Description class="dialog-description">
              Enter a name for your new folder
            </Dialog.Description>

            <div class="mt-4">
              <label for="folder-name" class="label">
                Folder Name
              </label>
              <input
                id="folder-name"
                class="input input-b-std w-full"
                type="text"
                value={ctx.folderName()}
                onInput={(e) => ctx.setFolderName(e.currentTarget.value)}
                placeholder="Enter folder name"
                autofocus
              />
              {ctx.error() && <p class="input-t-error mt-1">{ctx.error()}</p>}
            </div>

            <div class="flex justify-end gap-2 mt-4">
              <button
                class="btn btn-outline"
                onClick={() => ctx.setIsFolderModalOpen(false)}
                disabled={ctx.isCreating()}
              >
                Cancel
              </button>
              <button
                class="btn btn-primary"
                onClick={() => ctx.handleCreateFolder(params.folderId)}
                disabled={ctx.isCreating()}
              >
                {ctx.isCreating() ? "Creating..." : "Create Folder"}
              </button>
            </div>
          </div>
        </Dialog.Content>
      </Dialog.Portal>
    </Dialog>
  );
}
