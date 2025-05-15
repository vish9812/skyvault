import { FolderContent } from "@sv/apis/media/models";
import { For } from "solid-js";
import GridItem from "./grid-item";

interface GridViewProps {
  content: FolderContent;
}

function GridView(props: GridViewProps) {
  return (
    <div class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 gap-4">
      {/* Folders */}
      <For each={props.content.folderPage.items}>
        {(folder) => <GridItem type="folder" name={folder.name} />}
      </For>

      {/* Files */}
      <For each={props.content.filePage.items}>
        {(file) => (
          <GridItem type="file" name={file.name} preview={file.preview} />
        )}
      </For>
    </div>
  );
}

export default GridView;
