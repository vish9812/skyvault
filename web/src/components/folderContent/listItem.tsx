import type { FileInfo, FolderInfo } from "@sv/apis/media/models";
import { FileIcon } from "@sv/components/icons";
import { FOLDER_CONTENT_TYPES } from "@sv/utils/consts";
import Format from "@sv/utils/format";
import useCtx from "./ctxProvider";

interface ListItemProps {
  type: typeof FOLDER_CONTENT_TYPES.FILE | typeof FOLDER_CONTENT_TYPES.FOLDER;
  item: FileInfo & FolderInfo;
}

function ListItem(props: ListItemProps) {
  const ctx = useCtx();

  const handleClick = () => {
    ctx.handleTap(props.type, props.item.id);
  };

  return (
    <div
      class="grid grid-cols-[2rem_1fr_6rem] md:grid-cols-[2rem_1fr_6rem_9rem] items-center py-3 px-3 hover:bg-bg-muted border-t border-border first:border-t-0 cursor-pointer"
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
    </div>
  );
}

export default ListItem;
