import type { JSX } from "solid-js";

export interface FileListItemProps {
  name: string;
  size: number;
  updatedAt: string;
  children?: JSX.Element;
}

function formatSize(size: number) {
  if (size < 1024) return `${size} B`;
  if (size < 1024 * 1024) return `${(size / 1024).toFixed(1)} KB`;
  if (size < 1024 * 1024 * 1024)
    return `${(size / (1024 * 1024)).toFixed(1)} MB`;
  return `${(size / (1024 * 1024 * 1024)).toFixed(2)} GB`;
}

function formatDate(date: string) {
  return new Date(date).toLocaleDateString();
}

const FileListItem = (props: FileListItemProps): JSX.Element => {
  return (
    <div class="flex items-center py-3 px-4 hover:bg-gray-50 group relative">
      <span class="material-symbols-outlined text-gray-400 mr-3 text-2xl">
        description
      </span>
      <div class="flex-1 min-w-0">
        <div class="text-base font-semibold text-gray-900 truncate">
          {props.name}
        </div>
        <div class="text-xs text-gray-400 mt-0.5">
          {formatSize(props.size)} • {formatDate(props.updatedAt)}
        </div>
      </div>
      {props.children}
    </div>
  );
};

export default FileListItem;
