import { Show } from "solid-js";
import { FolderContent as FolderContentType } from "@sv/apis/media/models";
import EmptyState from "./emptyState";
import GridSkeleton from "./gridSkeleton";
import ListSkeleton from "./listSkeleton";
import ListView from "./listView";
import GridView from "./gridView";
import { CtxProvider } from "./ctxProvider";

interface FolderContentProps {
  content: FolderContentType | undefined;
  isListView: boolean;
  loading: boolean;
}

function FolderContentWithCtx(props: FolderContentProps) {
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

function FolderContent(props: FolderContentProps) {
  return (
    <CtxProvider>
      <FolderContentWithCtx {...props} />
    </CtxProvider>
  );
}

export default FolderContent;
