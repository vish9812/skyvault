import { FOLDER_CONTENT_TYPES } from "@sv/utils/consts";
import { formatDate, formatFileSize } from "@sv/utils/format";
import { getFileIcon } from "@sv/utils/icons";
import useViewModel from "./useViewModel";
import type { FileInfo, FolderInfo } from "@sv/apis/media/models";

interface ListItemProps {
  type: typeof FOLDER_CONTENT_TYPES.FILE | typeof FOLDER_CONTENT_TYPES.FOLDER;
  item: FileInfo & FolderInfo;
}

function ListItem(props: ListItemProps) {
  const { handleFolderNavigation } = useViewModel();

  const handleDoubleClick = () => {
    if (props.type === FOLDER_CONTENT_TYPES.FOLDER) {
      handleFolderNavigation(props.item.id);
    }
  };

  return (
    <div
      class="grid grid-cols-[2rem_1fr_6rem] md:grid-cols-[2rem_1fr_6rem_9rem] items-center py-3 px-3 hover:bg-bg-muted border-t border-border first:border-t-0 cursor-pointer"
      onDblClick={handleDoubleClick}
    >
      <span>
        {getFileIcon(
          props.type === FOLDER_CONTENT_TYPES.FOLDER,
          props.item.category
        )}
      </span>
      <div class="text-left font-medium text-neutral truncate">
        {props.item.name}
      </div>
      <div>
        {props.type === FOLDER_CONTENT_TYPES.FOLDER
          ? "-"
          : formatFileSize(props.item.size)}
      </div>
      <div class="hidden md:block">{formatDate(props.item.updatedAt)}</div>
    </div>
  );
}

export default ListItem;
