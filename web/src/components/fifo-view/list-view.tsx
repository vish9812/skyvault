import { FolderContent } from "@sv/apis/media/models";
import { For } from "solid-js";
import ListItem from "./list-item";

interface ListViewProps {
  content: FolderContent;
}

function ListView(props: ListViewProps) {
  return (
    <div class="flex flex-col border border-border rounded-lg overflow-hidden">
      {/* Header */}
      <div class="flex items-center py-3 px-3 bg-bg-muted border-b border-border font-medium text-sm text-neutral-light">
        <div class="w-6 mr-3"></div>
        <div class="flex-1">Name</div>
        <div class="w-24 text-right">Size</div>
        <div class="w-36 text-right">Last modified</div>
        <div class="w-8"></div>
      </div>

      {/* Folder list items */}
      <For each={props.content.folderPage.items}>
        {(folder) => (
          <ListItem
            type="folder"
            name={folder.name}
            updatedAt={folder.updatedAt}
          />
        )}
      </For>

      {/* File list items */}
      <For each={props.content.filePage.items}>
        {(file) => (
          <ListItem
            type="file"
            name={file.name}
            size={file.size}
            fileCategory={file.category}
            updatedAt={file.updatedAt}
          />
        )}
      </For>
    </div>
  );
}

export default ListView;
