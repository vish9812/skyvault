import { FILE_CATEGORIES } from "@sv/utils/consts";
import icons from "@sv/utils/icons";

interface ListItemProps {
  type: "file" | "folder";
  name: string;
  size?: number;
  fileCategory?: string;
  updatedAt?: string;
}

function ListItem(props: ListItemProps) {
  // Format file size to human-readable format
  const formatFileSize = (bytes?: number) => {
    if (bytes === undefined) return "-";
    if (bytes === 0) return "0 B";

    const units = ["B", "KB", "MB", "GB", "TB"];
    const i = Math.floor(Math.log(bytes) / Math.log(1024));
    return `${(bytes / Math.pow(1024, i)).toFixed(1)} ${units[i]}`;
  };

  // Format date to readable format
  const formatDate = (dateString?: string) => {
    if (!dateString) return "-";
    const date = new Date(dateString);
    return date.toLocaleDateString("en-US", {
      year: "numeric",
      month: "short",
      day: "numeric",
    });
  };

  const getFileIcon = () => {
    if (props.type === "folder") {
      return icons.folder({
        size: 5,
        color: "text-primary",
      });
    }

    switch (props.fileCategory) {
      case FILE_CATEGORIES.IMAGES:
        return icons.image({
          size: 5,
          color: "text-neutral-light",
        });
      case FILE_CATEGORIES.VIDEOS:
        return icons.video({
          size: 5,
          color: "text-neutral-light",
        });
      case FILE_CATEGORIES.AUDIOS:
        return icons.audio({
          size: 5,
          color: "text-neutral-light",
        });
      case FILE_CATEGORIES.DOCUMENTS:
        return icons.document({
          size: 5,
          color: "text-neutral-light",
        });
      default:
        return icons.file({
          size: 5,
          color: "text-neutral-light",
        });
    }
  };

  return (
    <div class="flex items-center py-3 px-4 hover:bg-bg-muted group border-t border-border first:border-t-0">
      <span
        class={`mr-3 ${
          props.type === "folder" ? "text-primary" : "text-neutral-light"
        }`}
      >
        {getFileIcon()}
      </span>
      <div class="flex-1 min-w-0">
        <div class="text-base font-medium text-neutral truncate">
          {props.name}
        </div>
      </div>
      <div class="w-24 text-sm text-neutral-light text-right">
        {props.type === "folder" ? "-" : formatFileSize(props.size)}
      </div>
      <div class="w-36 text-sm text-neutral-light text-right">
        {formatDate(props.updatedAt)}
      </div>
      <div class="w-8 opacity-0 group-hover:opacity-100 transition-opacity">
        <span class="material-symbols-outlined text-neutral-light cursor-pointer hover:text-primary">
          more_vert
        </span>
      </div>
    </div>
  );
}

export default ListItem;
