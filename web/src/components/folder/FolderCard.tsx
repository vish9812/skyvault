import type { JSX } from "solid-js";

export interface FolderCardProps {
  name: string;
  children?: JSX.Element;
}

const FolderCard = (props: FolderCardProps): JSX.Element => {
  return (
    <div class="bg-white rounded-lg shadow-sm border border-gray-200 p-4 hover:shadow-md transition-shadow relative group flex flex-col items-center">
      <span class="material-symbols-outlined text-4xl text-primary mb-2">
        folder
      </span>
      <div class="font-semibold text-base text-center text-gray-900 truncate w-full">
        {props.name}
      </div>
      <div class="text-xs text-gray-400 mt-1 mb-1">Folder</div>
      <div class="absolute top-2 right-2 opacity-80 group-hover:opacity-100 transition-opacity">
        {props.children}
      </div>
    </div>
  );
};

export default FolderCard;
