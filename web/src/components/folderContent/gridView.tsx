import type {
  FileInfo,
  FolderInfo,
  FolderContent,
} from "@sv/apis/media/models";
import { FOLDER_CONTENT_TYPES } from "@sv/utils/consts";
import { For } from "solid-js";
import GridItem from "./gridItem";

interface GridViewProps {
  content: FolderContent;
}

function GridView(props: GridViewProps) {
  return (
    <div class="p-4 flex flex-wrap gap-4 text-sm text-neutral-light">
      {/* Folders */}
      <For each={props.content.folderPage.items}>
        {(folder) => (
          <GridItem
            type={FOLDER_CONTENT_TYPES.FOLDER}
            item={folder as FileInfo & FolderInfo}
          />
        )}
      </For>

      {/* Files */}
      <For each={props.content.filePage.items}>
        {(file) => (
          <GridItem
            type={FOLDER_CONTENT_TYPES.FILE}
            item={file as FileInfo & FolderInfo}
          />
        )}
      </For>
    </div>
  );
}

export default GridView;
