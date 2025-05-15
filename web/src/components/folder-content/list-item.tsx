import { FOLDER_CONTENT_TYPES } from "@sv/utils/consts";
import { formatDate, formatFileSize } from "@sv/utils/format";
import { getFileIcon } from "@sv/utils/icons";

interface ListItemProps {
  type: typeof FOLDER_CONTENT_TYPES.FILE | typeof FOLDER_CONTENT_TYPES.FOLDER;
  name: string;
  size?: number;
  fileCategory?: string;
  updatedAt?: string;
}

function ListItem(props: ListItemProps) {
  return (
    <div class="grid grid-cols-[2rem_1fr_6rem] md:grid-cols-[2rem_1fr_6rem_9rem] items-center py-3 px-3 hover:bg-bg-muted border-t border-border first:border-t-0 cursor-pointer">
      <span>
        {getFileIcon(
          props.type === FOLDER_CONTENT_TYPES.FOLDER,
          props.fileCategory
        )}
      </span>
      <div class="text-left font-medium text-neutral truncate">
        {props.name}
      </div>
      <div>
        {props.type === FOLDER_CONTENT_TYPES.FOLDER
          ? "-"
          : formatFileSize(props.size)}
      </div>
      <div class="hidden md:block">{formatDate(props.updatedAt)}</div>
    </div>
  );
}

export default ListItem;
