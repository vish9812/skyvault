import { Button } from "@kobalte/core/button";
import { trashFiles, trashFolders } from "@sv/apis/media";
import Dialog from "@sv/components/ui/dialog";
import { FOLDER_CONTENT_TYPES } from "@sv/utils/consts";
import { defaultErrorMessage } from "@sv/utils/errors";
import Format from "@sv/utils/format";
import { createEffect, createSignal, Show } from "solid-js";
import type { DialogProps } from "./types";

function TrashDialog(props: DialogProps) {
  const [error, setError] = createSignal("");
  const [isLoading, setIsLoading] = createSignal(false);

  createEffect(() => {
    if (props.open) {
      setError("");
      setIsLoading(false);
    }
  });

  const handleTrash = async () => {
    if (isLoading()) return;

    setIsLoading(true);
    setError("");
    try {
      if (props.type === FOLDER_CONTENT_TYPES.FILE) {
        await trashFiles([props.item.id]);
      } else {
        await trashFolders([props.item.id]);
      }
      props.onDone();
      props.onClose();
    } catch (err) {
      if (err instanceof Error) {
        setError(defaultErrorMessage(err.message));
      } else {
        setError(`Failed to move the ${props.itemLabel} to trash.`);
      }
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <Dialog
      open={props.open}
      onClose={props.onClose}
      title={`Move ${Format.capitalize(props.itemLabel)} to Trash`}
      description={`"${props.item.name}" will be moved to trash.`}
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
              "btn-disabled": isLoading(),
              "btn-error": !isLoading(),
            }}
            onClick={handleTrash}
            disabled={isLoading()}
          >
            {isLoading() ? "Moving..." : "Move to Trash"}
          </Button>
        </>
      }
    >
      <div class="space-y-3">
        <p class="text-sm text-neutral">
          You can restore this {props.itemLabel} later from Trash.
        </p>
        <Show when={error()}>
          <p class="input-t-error">{error()}</p>
        </Show>
      </div>
    </Dialog>
  );
}

export default TrashDialog;
