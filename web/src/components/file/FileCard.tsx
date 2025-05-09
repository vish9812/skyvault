import type { JSX } from "solid-js";

export interface FileCardProps {
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

const FileCard = (props: FileCardProps): JSX.Element => {
  return (
    <div class="bg-white rounded-lg shadow-sm border border-gray-200 p-4 hover:shadow-md transition-shadow relative group flex flex-col items-center">
      <span class="material-symbols-outlined text-4xl text-gray-400 mb-2">
        description
      </span>
      <div class="font-semibold text-base text-center text-gray-900 truncate w-full">
        {props.name}
      </div>
      <div class="text-xs text-gray-400 mt-1 mb-1">
        {formatSize(props.size)} • {formatDate(props.updatedAt)}
      </div>
      <div class="absolute top-2 right-2 opacity-80 group-hover:opacity-100 transition-opacity">
        {props.children}
      </div>
    </div>
  );
};

export default FileCard;
