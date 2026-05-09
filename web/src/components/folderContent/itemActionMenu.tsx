import { Button } from "@kobalte/core/button";
import { DropdownMenu } from "@kobalte/core/dropdown-menu";
import { TextField } from "@kobalte/core/text-field";
import {
  downloadFile,
  fetchFolderContent,
  fetchFolderInfo,
  moveFile,
  moveFolder,
  renameFile,
  renameFolder,
  trashFiles,
  trashFolders,
} from "@sv/apis/media";
import type { FileInfo, FolderInfo } from "@sv/apis/media/models";
import Icon from "@sv/components/icons";
import Dialog from "@sv/components/ui/dialog";
import useAppCtx from "@sv/store/appCtxProvider";
import {
  FOLDER_CONTENT_TYPES,
  ROOT_FOLDER_ID,
  ROOT_FOLDER_NAME,
} from "@sv/utils/consts";
import { COMMON_ERR_KEYS, defaultErrorMessage } from "@sv/utils/errors";
import Validate, { VALIDATIONS } from "@sv/utils/validate";
import { createEffect, createSignal, For, Show } from "solid-js";
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
  const [isMoveOpen, setIsMoveOpen] = createSignal(false);
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
          onPointerDown={(e) => e.stopPropagation()}
          onClick={(e) => e.stopPropagation()}
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
            <DropdownMenu.Item
              class="dropdown-item"
              onSelect={() => setIsMoveOpen(true)}
            >
              <span class="flex items-center gap-2">
                <Icon name="move" size={5} color="text-neutral-light" />
                Move
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
      <MoveDialog
        type={props.type}
        item={props.item}
        itemLabel={itemLabel()}
        open={isMoveOpen()}
        onClose={() => setIsMoveOpen(false)}
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

function MoveDialog(props: DialogProps) {
  const [currentFolderId, setCurrentFolderId] = createSignal(ROOT_FOLDER_ID);
  const [folderInfo, setFolderInfo] = createSignal<FolderInfo | null>(null);
  const [childFolders, setChildFolders] = createSignal<FolderInfo[]>([]);
  const [error, setError] = createSignal("");
  const [isBrowsing, setIsBrowsing] = createSignal(false);
  const [isMoving, setIsMoving] = createSignal(false);
  let loadSeq = 0;

  const currentParentId = () => {
    if (props.type === FOLDER_CONTENT_TYPES.FILE) {
      return props.item.folderId ?? ROOT_FOLDER_ID;
    }
    return props.item.parentFolderId ?? ROOT_FOLDER_ID;
  };

  const isSameParent = () => currentFolderId() === currentParentId();
  const isMoveDisabled = () => isSameParent() || isBrowsing() || isMoving();

  const breadcrumbs = () => {
    const info = folderInfo();
    if (!info || info.id === ROOT_FOLDER_ID) {
      return [{ id: ROOT_FOLDER_ID, name: ROOT_FOLDER_NAME }];
    }

    return [
      { id: ROOT_FOLDER_ID, name: ROOT_FOLDER_NAME },
      ...info.ancestors.toReversed(),
      { id: info.id, name: info.name },
    ];
  };

  createEffect(() => {
    if (props.open) {
      setCurrentFolderId(currentParentId());
      setError("");
      setIsBrowsing(false);
      setIsMoving(false);
    }
  });

  createEffect(() => {
    if (!props.open) return;

    const seq = ++loadSeq;
    const id = currentFolderId();
    setIsBrowsing(true);
    setError("");

    Promise.all([fetchFolderInfo(id), fetchFolderContent(id)])
      .then(([info, content]) => {
        if (seq !== loadSeq) return;
        setFolderInfo(info);
        setChildFolders(content.folderPage.items);
      })
      .catch((err) => {
        if (seq !== loadSeq) return;
        if (err instanceof Error) {
          setError(defaultErrorMessage(err.message));
        } else {
          setError("Failed to load folders.");
        }
        setChildFolders([]);
      })
      .finally(() => {
        if (seq === loadSeq) {
          setIsBrowsing(false);
        }
      });
  });

  const canOpenFolder = (folder: FolderInfo) =>
    props.type !== FOLDER_CONTENT_TYPES.FOLDER || folder.id !== props.item.id;

  const handleMove = async () => {
    if (isMoveDisabled()) return;

    setIsMoving(true);
    setError("");
    try {
      const destinationFolderId =
        currentFolderId() === ROOT_FOLDER_ID ? "" : currentFolderId();

      if (props.type === FOLDER_CONTENT_TYPES.FILE) {
        await moveFile(props.item.id, destinationFolderId);
      } else {
        await moveFolder(props.item.id, destinationFolderId);
      }

      props.onDone();
      props.onClose();
    } catch (err) {
      if (err instanceof Error && err.message === COMMON_ERR_KEYS.INVALID) {
        setError(
          "Choose a different destination. Folders cannot be moved into themselves or their descendants."
        );
      } else if (err instanceof Error) {
        setError(defaultErrorMessage(err.message));
      } else {
        setError(`Failed to move the ${props.itemLabel}.`);
      }
    } finally {
      setIsMoving(false);
    }
  };

  return (
    <Dialog
      open={props.open}
      onClose={props.onClose}
      title={`Move ${capitalize(props.itemLabel)}`}
      description={`Choose a destination for "${props.item.name}".`}
      size="lg"
      actions={
        <>
          <Button
            class="btn btn-outline"
            onClick={props.onClose}
            disabled={isMoving()}
          >
            Cancel
          </Button>
          <Button
            classList={{
              btn: true,
              "btn-disabled": isMoveDisabled(),
              "btn-primary": !isMoveDisabled(),
            }}
            onClick={handleMove}
            disabled={isMoveDisabled()}
          >
            {isMoving() ? "Moving..." : "Move here"}
          </Button>
        </>
      }
    >
      <div class="space-y-4">
        <div class="flex flex-wrap items-center gap-1 text-sm">
          <For each={breadcrumbs()}>
            {(crumb, index) => (
              <>
                <Show when={index() > 0}>
                  <span class="text-neutral-lighter">/</span>
                </Show>
                <button
                  type="button"
                  classList={{
                    link: crumb.id !== currentFolderId(),
                    "font-semibold text-neutral":
                      crumb.id === currentFolderId(),
                  }}
                  onClick={() => setCurrentFolderId(crumb.id)}
                  disabled={crumb.id === currentFolderId() || isBrowsing()}
                >
                  {crumb.name}
                </button>
              </>
            )}
          </For>
        </div>

        <div class="min-h-48 rounded-md border border-border-strong overflow-hidden">
          <Show
            when={!isBrowsing()}
            fallback={
              <div class="flex-center min-h-48 text-sm text-neutral-lighter">
                Loading folders...
              </div>
            }
          >
            <Show
              when={childFolders().length > 0}
              fallback={
                <div class="flex-center min-h-48 text-sm text-neutral-lighter">
                  This destination has no child folders.
                </div>
              }
            >
              <div class="divide-y divide-border">
                <For each={childFolders()}>
                  {(folder) => {
                    const disabled = () => !canOpenFolder(folder);

                    return (
                      <button
                        type="button"
                        classList={{
                          "w-full flex items-center justify-between gap-3 px-4 py-3 text-left transition-colors": true,
                          "hover:bg-bg-muted": !disabled(),
                          "cursor-not-allowed bg-bg-subtle text-neutral-lighter":
                            disabled(),
                        }}
                        disabled={disabled() || isBrowsing()}
                        onClick={() => setCurrentFolderId(folder.id)}
                      >
                        <span class="min-w-0 flex items-center gap-3">
                          <Icon
                            name="folder"
                            size={5}
                            color={
                              disabled()
                                ? "text-neutral-lighter"
                                : "text-primary"
                            }
                          />
                          <span class="truncate">{folder.name}</span>
                        </span>
                        <Show when={disabled()}>
                          <span class="shrink-0 text-xs">Cannot select</span>
                        </Show>
                      </button>
                    );
                  }}
                </For>
              </div>
            </Show>
          </Show>
        </div>

        <Show when={isSameParent()}>
          <p class="text-xs text-neutral-lighter">
            This is the current location for the {props.itemLabel}.
          </p>
        </Show>
        <Show when={error()}>
          <p class="input-t-error">{error()}</p>
        </Show>
      </div>
    </Dialog>
  );
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
