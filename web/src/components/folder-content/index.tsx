import { Show } from "solid-js";
import { FolderContent as FolderContentType } from "@sv/apis/media/models";
import EmptyState from "./empty-state";
import GridSkeleton from "./grid-skeleton";
import ListSkeleton from "./list-skeleton";
import ListView from "./list-view";
import GridView from "./grid-view";

interface FolderContentProps {
  content: FolderContentType | undefined;
  isListView: boolean;
  loading: boolean;
}

function FolderContent(props: FolderContentProps) {
  return (
    <Show
      when={!props.loading}
      fallback={
        <Show when={props.isListView} fallback={<GridSkeleton />}>
          <ListSkeleton />
        </Show>
      }
    >
      <Show
        when={
          (props.content?.folderPage?.items?.length ?? 0) > 0 ||
          (props.content?.filePage?.items?.length ?? 0) > 0
        }
        fallback={<EmptyState />}
      >
        <Show
          when={props.isListView}
          fallback={<GridView content={props.content!} />}
        >
          <ListView content={props.content!} />
        </Show>
      </Show>
    </Show>
  );
}

export default FolderContent;
