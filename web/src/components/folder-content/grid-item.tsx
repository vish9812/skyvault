import { FOLDER_CONTENT_TYPES } from "@sv/utils/consts";
import { formatFileSize } from "@sv/utils/format";
import { getFileIcon } from "@sv/utils/icons";
import { Show } from "solid-js";

const iconSize = 10;

interface GridItemProps {
  type: typeof FOLDER_CONTENT_TYPES.FILE | typeof FOLDER_CONTENT_TYPES.FOLDER;
  name: string;
  preview?: string;
  size?: number;
  fileCategory?: string;
}

function GridItem(props: GridItemProps) {
  const formattedSize = formatFileSize(props.size);

  return (
    <div class="w-40 h-40 md:w-48 md:h-48 bg-white rounded-lg border border-border shadow-sm hover:border-primary hover:shadow-md transition-all cursor-pointer">
      {/* File/folder icon or preview */}
      <div class="flex-center h-28 md:h-34 rounded-t-lg border-b border-border bg-bg-subtle">
        {props.preview ? (
          <img
            src={`data:image/png;base64,${props.preview}`}
            alt={props.name}
            class="object-cover h-full w-full"
          />
        ) : (
          <span>
            {getFileIcon(
              props.type === FOLDER_CONTENT_TYPES.FOLDER,
              props.fileCategory,
              iconSize
            )}
          </span>
        )}
      </div>

      {/* File/folder info */}
      <div class="p-2 flex flex-col">
        <div class="font-medium text-neutral truncate">{props.name}</div>
        <Show when={formattedSize !== "-"}>
          <div class="caption">{formattedSize}</div>
        </Show>
      </div>
    </div>
  );
}

export default GridItem;
