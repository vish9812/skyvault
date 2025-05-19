import { Dialog } from "@kobalte/core/dialog";
import { TextField } from "@kobalte/core/text-field";
import useCtx from "./ctxProvider";
import { useParams } from "@solidjs/router";
import { Button } from "@kobalte/core/button";
import { createEffect } from "solid-js";

function CreateFolder() {
  const ctx = useCtx();
  const params = useParams();

  const isDisabled = () => !!ctx.error() || ctx.isCreating();

  createEffect(() => {
    console.log("isDisabled: ", isDisabled());
  });

  return (
    <Dialog
      open={ctx.isCreateFolderModalOpen()}
      onOpenChange={(isOpen) =>
        !isOpen && ctx.setIsCreateFolderModalOpen(false)
      }
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
              <TextField
                value={ctx.createFolderName()}
                onChange={ctx.handleCreateFolderNameChange}
                validationState={ctx.error() ? "invalid" : "valid"}
              >
                <TextField.Label class="label">Folder Name</TextField.Label>
                <TextField.Input
                  class="input"
                  classList={{
                    "input-b-std": !ctx.error(),
                    "input-b-error": !!ctx.error(),
                  }}
                  type="text"
                  placeholder="Enter folder name"
                  autocomplete="off"
                  autofocus
                />
                <TextField.ErrorMessage class="input-t-error">
                  {ctx.error()}
                </TextField.ErrorMessage>
              </TextField>
            </div>

            <div class="flex justify-end gap-2 mt-4">
              <Button
                class="btn btn-outline"
                onClick={() => ctx.setIsCreateFolderModalOpen(false)}
                disabled={ctx.isCreating()}
              >
                Cancel
              </Button>
              <Button
                class="btn"
                classList={{
                  "btn-disabled": isDisabled(),
                  "btn-primary": !isDisabled(),
                }}
                onClick={() => ctx.handleCreateFolder(params.folderId)}
                disabled={isDisabled()}
              >
                {ctx.isCreating() ? "Creating..." : "Create Folder"}
              </Button>
            </div>
          </div>
        </Dialog.Content>
      </Dialog.Portal>
    </Dialog>
  );
}

export default CreateFolder;
