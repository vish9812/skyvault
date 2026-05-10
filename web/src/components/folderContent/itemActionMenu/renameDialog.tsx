import { Button } from "@kobalte/core/button";
import { TextField } from "@kobalte/core/text-field";
import { renameFile, renameFolder } from "@sv/apis/media";
import Dialog from "@sv/components/ui/dialog";
import { FOLDER_CONTENT_TYPES } from "@sv/utils/consts";
import { COMMON_ERR_KEYS, defaultErrorMessage } from "@sv/utils/errors";
import Format from "@sv/utils/format";
import Validate, { VALIDATIONS } from "@sv/utils/validate";
import { createEffect, createSignal } from "solid-js";
import type { DialogProps } from "./types";

function RenameDialog(props: DialogProps) {
  const [name, setName] = createSignal(props.item.name);
  const [error, setError] = createSignal("");
  const [isLoading, setIsLoading] = createSignal(false);
  let inputRef!: HTMLInputElement;

  const isInvalidName = () => !Validate.name(name());
  const isUnchanged = () => name().trim() === props.item.name;
  const isDisabled = () => isInvalidName() || isUnchanged() || isLoading();

  createEffect(() => {
    if (props.open) {
      setName(props.item.name);
      setError("");
      setIsLoading(false);
      setTimeout(() => {
        inputRef?.focus();
        inputRef?.select();
      }, 100);
    }
  });

  const handleNameChange = (value: string) => {
    setName(value);
    let msg = "";
    if (!Validate.name(value)) {
      msg =
        !value || value.trim().length === 0
          ? `${Format.capitalize(props.itemLabel)} name is required`
          : `${Format.capitalize(props.itemLabel)} name must be less than ${VALIDATIONS.MAX_LENGTH} characters`;
    }
    setError(msg);
  };

  const handleRename = async () => {
    if (isDisabled()) return;

    setIsLoading(true);
    setError("");
    try {
      if (props.type === FOLDER_CONTENT_TYPES.FILE) {
        await renameFile(props.item.id, name().trim());
      } else {
        await renameFolder(props.item.id, name().trim());
      }
      props.onDone();
      props.onClose();
    } catch (err) {
      if (err instanceof Error && err.message === COMMON_ERR_KEYS.DUPLICATE) {
        setError(`A ${props.itemLabel} with this name already exists here.`);
      } else if (err instanceof Error) {
        setError(defaultErrorMessage(err.message));
      } else {
        setError(`Failed to rename the ${props.itemLabel}.`);
      }
    } finally {
      setIsLoading(false);
    }
  };

  const handleKeyDown = (e: KeyboardEvent) => {
    if (e.key === "Enter" && !isDisabled()) {
      e.preventDefault();
      handleRename();
    }
  };

  return (
    <Dialog
      open={props.open}
      onClose={props.onClose}
      title={`Rename ${Format.capitalize(props.itemLabel)}`}
      description={`Enter a new name for "${props.item.name}".`}
      size="md"
      actions={
        <>
          <Button
            class="btn btn-outline"
            onClick={props.onClose}
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
            onClick={handleRename}
            disabled={isDisabled()}
          >
            {isLoading() ? "Renaming..." : "Rename"}
          </Button>
        </>
      }
    >
      <TextField
        value={name()}
        onChange={handleNameChange}
        validationState={error() ? "invalid" : "valid"}
      >
        <TextField.Label class="label">
          {Format.capitalize(props.itemLabel)} Name
        </TextField.Label>
        <TextField.Input
          ref={inputRef}
          classList={{
            input: true,
            "input-b-std": !error(),
            "input-b-error": !!error(),
          }}
          type="text"
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

export default RenameDialog;
