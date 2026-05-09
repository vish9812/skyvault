import { DropdownMenu } from "@kobalte/core/dropdown-menu";
import { downloadFile } from "@sv/apis/media";
import Icon from "@sv/components/icons";
import useAppCtx from "@sv/store/appCtxProvider";
import { FOLDER_CONTENT_TYPES } from "@sv/utils/consts";
import { createSignal, Show } from "solid-js";
import useCtx from "../ctxProvider";
import MoveDialog from "./moveDialog";
import RenameDialog from "./renameDialog";
import TrashDialog from "./trashDialog";
import type { ActionItem, ItemType } from "./types";

interface Props {
  type: ItemType;
  item: ActionItem;
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

  const handleDone = () => {
    ctx.clearSelection();
    appCtx.refreshFolderContent();
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
        onDone={handleDone}
      />
      <MoveDialog
        type={props.type}
        item={props.item}
        itemLabel={itemLabel()}
        open={isMoveOpen()}
        onClose={() => setIsMoveOpen(false)}
        onDone={handleDone}
      />
      <TrashDialog
        type={props.type}
        item={props.item}
        itemLabel={itemLabel()}
        open={isTrashOpen()}
        onClose={() => setIsTrashOpen(false)}
        onDone={() => {
          handleDone();
          appCtx.refreshStorageUsage();
        }}
      />
    </>
  );
}

export default ItemActionMenu;

