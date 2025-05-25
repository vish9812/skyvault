import { FOLDER_CONTENT_TYPES } from "@sv/utils/consts";
import { Show } from "solid-js";
import useCtx from "./ctxProvider";
import type { FileInfo, FolderInfo } from "@sv/apis/media/models";
import format from "@sv/utils/format";
import { FileIcon } from "@sv/utils/icons";

interface GridItemProps {
  type: typeof FOLDER_CONTENT_TYPES.FILE | typeof FOLDER_CONTENT_TYPES.FOLDER;
  item: FileInfo & FolderInfo;
}

function GridItem(props: GridItemProps) {
  const ctx = useCtx();

  const handleClick = () => {
    ctx.handleTap(props.type, props.item.id);
  };

  return (
    <div
      class="w-40 h-40 md:w-48 md:h-48 bg-white rounded-lg border border-border shadow-sm hover:border-primary hover:shadow-md transition-all cursor-pointer"
      onClick={handleClick}
    >
      {/* File/folder icon or preview */}
      <div class="flex-center h-28 md:h-34 rounded-t-lg border-b border-border bg-bg-subtle">
        {props.item.preview ? (
          <img
            src={`data:image/png;base64,${props.item.preview}`}
            alt={props.item.name}
            class="object-cover h-full w-full"
          />
        ) : (
          <span>
            <FileIcon
              fileCategory={props.item.category}
              isFolder={props.type === FOLDER_CONTENT_TYPES.FOLDER}
              size={10}
            />
          </span>
        )}
      </div>

      {/* File/folder info */}
      <div class="p-2 flex flex-col">
        <div class="font-medium text-neutral truncate">{props.item.name}</div>
        <Show when={props.type === FOLDER_CONTENT_TYPES.FILE}>
          <div class="text-xs">{format.size(props.item.size)}</div>
        </Show>
      </div>
    </div>
  );
}

export default GridItem;
