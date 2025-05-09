import type { JSX } from "solid-js";

export interface FolderListItemProps {
  name: string;
  children?: JSX.Element;
}

const FolderListItem = (props: FolderListItemProps): JSX.Element => {
  return (
    <div class="flex items-center py-3 px-4 hover:bg-gray-50 group relative">
      <span class="material-symbols-outlined text-primary mr-3 text-2xl">
        folder
      </span>
      <div class="flex-1 min-w-0">
        <div class="text-base font-semibold text-gray-900 truncate">
          {props.name}
        </div>
        <div class="text-xs text-gray-400 mt-0.5">Folder</div>
      </div>
      {props.children}
    </div>
  );
};

export default FolderListItem;
