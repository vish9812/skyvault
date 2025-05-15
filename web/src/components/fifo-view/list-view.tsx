import { FolderContent } from "@sv/apis/media/models";
import { For } from "solid-js";
import ListItem from "./list-item";

interface ListViewProps {
  content: FolderContent;
}

function ListView(props: ListViewProps) {
  return (
    <div class="text-neutral-light text-sm text-right border border-border rounded-lg overflow-hidden max-w-screen-xl mx-auto">
      {/* Header */}
      <div class="grid grid-cols-[2rem_1fr_6rem] md:grid-cols-[2rem_1fr_6rem_9rem] items-center py-3 px-3 bg-bg-muted border-b border-border font-medium">
        <div></div>
        <div class="text-left">Name</div>
        <div>Size</div>
        <div class="hidden md:block">Last modified</div>
        <div></div>
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
