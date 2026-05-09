import { Button } from "@kobalte/core/button";
import { DropdownMenu } from "@kobalte/core/dropdown-menu";
import { TextField } from "@kobalte/core/text-field";
import {
  downloadFile,
  renameFile,
  renameFolder,
  trashFiles,
  trashFolders,
} from "@sv/apis/media";
import type { FileInfo, FolderInfo } from "@sv/apis/media/models";
import Icon from "@sv/components/icons";
import Dialog from "@sv/components/ui/dialog";
import useAppCtx from "@sv/store/appCtxProvider";
import { FOLDER_CONTENT_TYPES } from "@sv/utils/consts";
import { COMMON_ERR_KEYS, defaultErrorMessage } from "@sv/utils/errors";
import Validate, { VALIDATIONS } from "@sv/utils/validate";
import { createEffect, createSignal, Show } from "solid-js";
import useCtx from "./ctxProvider";

type ItemType =
  | typeof FOLDER_CONTENT_TYPES.FILE
  | typeof FOLDER_CONTENT_TYPES.FOLDER;

interface Props {
  type: ItemType;
  item: FileInfo & FolderInfo;
}

function ItemActionMenu(props: Props) {
  const appCtx = useAppCtx();
  const ctx = useCtx();
  const [isRenameOpen, setIsRenameOpen] = createSignal(false);
  const [isTrashOpen, setIsTrashOpen] = createSignal(false);
  const [isDownloading, setIsDownloading] = createSignal(false);

  const itemLabel = () =>
    props.type === FOLDER_CONTENT_TYPES.FILE ? "file" : "folder";

  const handleDownload = async () => {
    if (isDownloading() || props.type !== FOLDER_CONTENT_TYPES.FILE) return;

    setIsDownloading(true);
    try {
      await downloadFile(props.item.id, props.item.name);
      ctx.clearSelection();
    } catch (err) {
      console.error("Download failed:", err);
    } finally {
      setIsDownloading(false);
    }
  };

  return (
    <>
      <DropdownMenu placement="bottom-end">
        <DropdownMenu.Trigger
          class="w-8 h-8 rounded-md flex-center hover:bg-bg-muted transition-colors"
          title={`${props.item.name} actions`}
          onClick={(e) => {
            e.stopPropagation();
            ctx.handleTap({
              id: props.item.id,
              type: props.type,
              name: props.item.name,
            });
          }}
        >
          <Icon name="moreOptions" size={5} color="text-neutral-light" />
        </DropdownMenu.Trigger>
        <DropdownMenu.Portal>
          <DropdownMenu.Content
            class="bg-white rounded-lg shadow-md border border-border-strong min-w-[180px] py-2 z-20"
            onClick={(e) => e.stopPropagation()}
          >
            <Show when={props.type === FOLDER_CONTENT_TYPES.FILE}>
              <DropdownMenu.Item
                class="dropdown-item"
                disabled={isDownloading()}
                onSelect={handleDownload}
              >
                <span class="flex items-center gap-2">
                  <Icon name="download" size={5} color="text-neutral-light" />
                  {isDownloading() ? "Downloading..." : "Download"}
                </span>
              </DropdownMenu.Item>
            </Show>
            <DropdownMenu.Item
              class="dropdown-item"
              onSelect={() => setIsRenameOpen(true)}
            >
              <span class="flex items-center gap-2">
                <Icon name="rename" size={5} color="text-neutral-light" />
                Rename
              </span>
            </DropdownMenu.Item>
            <DropdownMenu.Separator class="border-border-strong my-2" />
            <DropdownMenu.Item
              class="dropdown-item text-error"
              onSelect={() => setIsTrashOpen(true)}
            >
              <span class="flex items-center gap-2">
                <Icon name="trash" size={5} color="text-error" />
                Move to trash
              </span>
            </DropdownMenu.Item>
          </DropdownMenu.Content>
        </DropdownMenu.Portal>
      </DropdownMenu>

      <RenameDialog
        type={props.type}
        item={props.item}
        itemLabel={itemLabel()}
        open={isRenameOpen()}
        onClose={() => setIsRenameOpen(false)}
        onDone={() => {
          ctx.clearSelection();
          appCtx.refreshFolderContent();
        }}
      />
      <TrashDialog
        type={props.type}
        item={props.item}
        itemLabel={itemLabel()}
        open={isTrashOpen()}
        onClose={() => setIsTrashOpen(false)}
        onDone={() => {
          ctx.clearSelection();
          appCtx.refreshFolderContent();
          appCtx.refreshStorageUsage();
        }}
      />
    </>
  );
}

interface DialogProps {
  type: ItemType;
  item: FileInfo & FolderInfo;
  itemLabel: string;
  open: boolean;
  onClose: () => void;
  onDone: () => void;
}

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
          ? `${capitalize(props.itemLabel)} name is required`
          : `${capitalize(props.itemLabel)} name must be less than ${VALIDATIONS.MAX_LENGTH} characters`;
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
      title={`Rename ${capitalize(props.itemLabel)}`}
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
          {capitalize(props.itemLabel)} Name
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
      title={`Move ${capitalize(props.itemLabel)} to Trash`}
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

function capitalize(value: string): string {
  return value.charAt(0).toUpperCase() + value.slice(1);
}

export default ItemActionMenu;
