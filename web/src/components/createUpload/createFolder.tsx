import { Button } from "@kobalte/core/button";
import { TextField } from "@kobalte/core/text-field";
import { createFolder } from "@sv/apis/media";
import Dialog from "@sv/components/ui/dialog";
import useAppCtx from "@sv/store/appCtxProvider";
import { COMMON_ERR_KEYS } from "@sv/utils/errors";
import Validate, { VALIDATIONS } from "@sv/utils/validate";
import { createEffect, createSignal } from "solid-js";

interface Props {
  isModalOpen: boolean;
  closeModal: () => void;
}

function CreateFolder(props: Props) {
  const appCtx = useAppCtx();

  const [name, setName] = createSignal("");
  const [isLoading, setIsLoading] = createSignal(false);
  const [error, setError] = createSignal("");
  let inputRef!: HTMLInputElement;

  const isInvalidName = () => !Validate.name(name());
  const isDisabled = () => isInvalidName() || isLoading();

  // Reset state when the modal is closed
  createEffect(() => {
    if (!props.isModalOpen) {
      setName("");
      setError("");
      setIsLoading(false);
    }
  });

  // Auto focus when modal opens
  createEffect(() => {
    if (props.isModalOpen && inputRef) {
      setTimeout(() => inputRef?.focus(), 100);
    }
  });

  const handleNameChange = (name: string) => {
    setName(name);
    let msg = "";
    if (isInvalidName()) {
      msg =
        !name || name.trim().length === 0
          ? "Folder name is required"
          : `Folder name must be less than ${VALIDATIONS.MAX_LENGTH} characters`;
    }
    setError(msg);
  };

  const handleCreateFolder = async () => {
    setIsLoading(true);
    if (isInvalidName()) {
      setIsLoading(false);
      return;
    }

    try {
      await createFolder(appCtx.currentFolderId(), name().trim());
      props.closeModal();
    } catch (err) {
      if (err instanceof Error && err.message === COMMON_ERR_KEYS.DUPLICATE) {
        setError("A folder with this name already exists in this folder.");
      } else {
        setError("Failed to create the folder");
      }
    } finally {
      setIsLoading(false);
    }
  };

  const handleKeyDown = (e: KeyboardEvent) => {
    if (e.key === "Enter" && !isDisabled()) {
      e.preventDefault();
      handleCreateFolder();
    }
  };

  return (
    <Dialog
      open={props.isModalOpen}
      onClose={props.closeModal}
      title="Create New Folder"
      description="Enter a name for your new folder"
      size="md"
      actions={
        <>
          <Button
            class="btn btn-outline"
            onClick={props.closeModal}
            disabled={isLoading()}
          >
            Close
          </Button>
          <Button
            classList={{
              btn: true,
              "btn-disabled": isDisabled(),
              "btn-primary": !isDisabled(),
            }}
            onClick={handleCreateFolder}
            disabled={isDisabled()}
          >
            {isLoading() ? "Creating..." : "Create Folder"}
          </Button>
        </>
      }
    >
      <TextField
        value={name()}
        onChange={handleNameChange}
        validationState={error() ? "invalid" : "valid"}
      >
        <TextField.Label class="label">Folder Name</TextField.Label>
        <TextField.Input
          ref={inputRef}
          classList={{
            input: true,
            "input-b-std": !error(),
            "input-b-error": !!error(),
          }}
          type="text"
          placeholder="Enter folder name"
          autocomplete="off"
          onKeyDown={handleKeyDown}
        />
        <TextField.ErrorMessage class="input-t-error">
          {error()}
        </TextField.ErrorMessage>
      </TextField>
    </Dialog>
  );
}

export default CreateFolder;
