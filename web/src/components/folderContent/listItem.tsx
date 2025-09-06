import type { FileInfo, FolderInfo } from "@sv/apis/media/models";
import { FileIcon } from "@sv/components/icons";
import Icon from "@sv/components/icons";
import { FOLDER_CONTENT_TYPES } from "@sv/utils/consts";
import Format from "@sv/utils/format";
import { Show, createSignal } from "solid-js";
import { downloadFile } from "@sv/apis/media";
import useCtx from "./ctxProvider";

interface ListItemProps {
  type: typeof FOLDER_CONTENT_TYPES.FILE | typeof FOLDER_CONTENT_TYPES.FOLDER;
  item: FileInfo & FolderInfo;
}

function ListItem(props: ListItemProps) {
  const ctx = useCtx();
  const [isDownloading, setIsDownloading] = createSignal(false);

  const handleClick = () => {
    ctx.handleTap(props.type, props.item.id);
  };

  const handleDownload = async (e: Event) => {
    e.stopPropagation(); // Prevent triggering the main click
    if (isDownloading() || props.type !== FOLDER_CONTENT_TYPES.FILE) return;
    
    setIsDownloading(true);
    try {
      await downloadFile(props.item.id, props.item.name);
    } catch (error) {
      console.error('Download failed:', error);
    } finally {
      setIsDownloading(false);
    }
  };

  return (
    <div
      class="grid grid-cols-[2rem_1fr_6rem_2rem] md:grid-cols-[2rem_1fr_6rem_9rem_2rem] items-center py-3 px-3 hover:bg-bg-muted border-t border-border first:border-t-0 cursor-pointer group"
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
        <Show when={props.type === FOLDER_CONTENT_TYPES.FILE}>
          <button
            class="w-6 h-6 rounded-full flex-center opacity-0 group-hover:opacity-100 transition-opacity hover:bg-bg-subtle"
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
