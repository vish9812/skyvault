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
