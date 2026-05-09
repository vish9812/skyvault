import type { FileInfo, FolderInfo } from "@sv/apis/media/models";
import { FileIcon } from "@sv/components/icons";
import { FOLDER_CONTENT_TYPES } from "@sv/utils/consts";
import Format from "@sv/utils/format";
import { Show } from "solid-js";
import ItemActionMenu from "./itemActionMenu";
import useCtx from "./ctxProvider";

interface GridItemProps {
  type: typeof FOLDER_CONTENT_TYPES.FILE | typeof FOLDER_CONTENT_TYPES.FOLDER;
  item: FileInfo & FolderInfo;
}

function GridItem(props: GridItemProps) {
  const ctx = useCtx();

  const isSelected = () => ctx.selectedItem()?.id === props.item.id;

  const handleClick = () => {
    ctx.handleTap({
      id: props.item.id,
      type: props.type,
      name: props.item.name,
    });
  };

  return (
    <div
      class={`w-40 h-40 md:w-48 md:h-48 bg-white rounded-lg border shadow-sm transition-all cursor-pointer relative ${
        isSelected()
          ? "border-primary shadow-md"
          : "border-border hover:border-primary hover:shadow-md"
      }`}
      onClick={handleClick}
    >
      <div class="absolute top-2 right-2 bg-white rounded-lg shadow-lg border border-border p-1 z-10">
        <ItemActionMenu type={props.type} item={props.item} />
      </div>

      {/* File/folder icon or preview */}
      <div class="flex-center h-28 md:h-34 rounded-t-lg border-b border-border bg-bg-subtle">
        {props.item.previewBase64 ? (
          <img
            src={`data:image/png;base64,${props.item.previewBase64}`}
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
          <div class="text-xs">{Format.size(props.item.size)}</div>
        </Show>
      </div>
    </div>
  );
}

export default GridItem;
