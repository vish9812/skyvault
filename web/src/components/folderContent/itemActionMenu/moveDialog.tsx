import { Button } from "@kobalte/core/button";
import {
  fetchChildFolders,
  fetchFolderInfo,
  moveFile,
  moveFolder,
} from "@sv/apis/media";
import type { FolderInfo } from "@sv/apis/media/models";
import Icon from "@sv/components/icons";
import Dialog from "@sv/components/ui/dialog";
import {
  FOLDER_CONTENT_TYPES,
  ROOT_FOLDER_ID,
  ROOT_FOLDER_NAME,
} from "@sv/utils/consts";
import { COMMON_ERR_KEYS, defaultErrorMessage } from "@sv/utils/errors";
import Format from "@sv/utils/format";
import { createEffect, createSignal, For, Show } from "solid-js";
import type { DialogProps } from "./types";

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

    Promise.all([fetchFolderInfo(id), fetchChildFolders(id)])
      .then(([info, folders]) => {
        if (seq !== loadSeq) return;
        setFolderInfo(info);
        setChildFolders(folders.items);
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
      title={`Move ${Format.capitalize(props.itemLabel)}`}
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

export default MoveDialog;
