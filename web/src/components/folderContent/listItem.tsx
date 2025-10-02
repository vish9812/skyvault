import { downloadFile } from "@sv/apis/media";
import type { FileInfo, FolderInfo } from "@sv/apis/media/models";
import Icon, { FileIcon } from "@sv/components/icons";
import { FOLDER_CONTENT_TYPES } from "@sv/utils/consts";
import Format from "@sv/utils/format";
import { Show, createSignal } from "solid-js";
import useCtx from "./ctxProvider";

interface ListItemProps {
  type: typeof FOLDER_CONTENT_TYPES.FILE | typeof FOLDER_CONTENT_TYPES.FOLDER;
  item: FileInfo & FolderInfo;
}

function ListItem(props: ListItemProps) {
  const ctx = useCtx();
  const [isDownloading, setIsDownloading] = createSignal(false);

  const isSelected = () => ctx.selectedItem()?.id === props.item.id;

  const handleClick = () => {
    ctx.handleTap({
      id: props.item.id,
      type: props.type,
      name: props.item.name,
    });
  };

  const handleDownload = async (e: Event) => {
    e.stopPropagation();
    if (isDownloading() || props.type !== FOLDER_CONTENT_TYPES.FILE) return;

    setIsDownloading(true);
    try {
      await downloadFile(props.item.id, props.item.name);
      ctx.clearSelection(); // Clear selection after download
    } catch (error) {
      console.error("Download failed:", error);
    } finally {
      setIsDownloading(false);
    }
  };

  return (
    <div
      class={`grid grid-cols-[2rem_1fr_6rem_2rem] md:grid-cols-[2rem_1fr_6rem_9rem_2rem] items-center py-3 px-3 border-t border-border first:border-t-0 cursor-pointer transition-colors ${
        isSelected() ? "bg-primary-lighter" : "hover:bg-bg-muted"
      }`}
      onClick={handleClick}
    >
      <span>
        <FileIcon
          fileCategory={props.item.category}
          isFolder={props.type === FOLDER_CONTENT_TYPES.FOLDER}
        />
      </span>
      <div class="text-left font-medium text-neutral truncate">
        {props.item.name}
      </div>
      <div>
        {props.type === FOLDER_CONTENT_TYPES.FOLDER
          ? "-"
          : Format.size(props.item.size)}
      </div>
      <div class="hidden md:block">{Format.date(props.item.updatedAt)}</div>
      <div class="flex justify-center">
        <Show when={isSelected() && props.type === FOLDER_CONTENT_TYPES.FILE}>
          <button
            class="w-8 h-8 rounded-md flex-center hover:bg-bg-muted transition-colors"
            onClick={handleDownload}
            disabled={isDownloading()}
            title="Download file"
          >
            <Icon
              name="download"
              size={4}
              color={isDownloading() ? "text-neutral-lighter" : "text-primary"}
            />
          </button>
        </Show>
      </div>
    </div>
  );
}

export default ListItem;
