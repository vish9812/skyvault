import { DropdownMenu } from "@kobalte/core/dropdown-menu";
import type { JSX } from "solid-js";

interface ItemContextMenuProps {
  type: "file" | "folder";
  id: number;
  name: string;
}

const ItemContextMenu = (props: ItemContextMenuProps): JSX.Element => {
  return (
    <DropdownMenu>
      <DropdownMenu.Trigger class="text-gray-400 hover:text-gray-600 rounded-full p-1 absolute top-2 right-2 opacity-0 group-hover:opacity-100 transition-opacity">
        <span class="material-symbols-outlined text-base">more_vert</span>
      </DropdownMenu.Trigger>
      <DropdownMenu.Content class="bg-white rounded-lg shadow-lg border border-gray-200 py-1 min-w-[180px]">
        <DropdownMenu.Item class="px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 cursor-pointer flex items-center gap-2">
          <span class="material-symbols-outlined text-gray-500 text-sm">
            drive_file_rename_outline
          </span>
          <span>Rename</span>
        </DropdownMenu.Item>
        <DropdownMenu.Item class="px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 cursor-pointer flex items-center gap-2">
          <span class="material-symbols-outlined text-gray-500 text-sm">
            drive_file_move
          </span>
          <span>Move</span>
        </DropdownMenu.Item>
        <DropdownMenu.Item class="px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 cursor-pointer flex items-center gap-2">
          <span class="material-symbols-outlined text-gray-500 text-sm">
            share
          </span>
          <span>Share</span>
        </DropdownMenu.Item>
        <DropdownMenu.Item class="px-4 py-2 text-sm text-red-600 hover:bg-gray-100 cursor-pointer flex items-center gap-2">
          <span class="material-symbols-outlined text-red-500 text-sm">
            delete
          </span>
          <span>Delete</span>
        </DropdownMenu.Item>
      </DropdownMenu.Content>
    </DropdownMenu>
  );
};

export default ItemContextMenu;
