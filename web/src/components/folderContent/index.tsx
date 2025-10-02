import { Show } from "solid-js";
import { FolderContent as FolderContentType } from "@sv/apis/media/models";
import EmptyState from "./emptyState";
import GridSkeleton from "./gridSkeleton";
import ListSkeleton from "./listSkeleton";
import ListView from "./listView";
import GridView from "./gridView";
import { CtxProvider } from "./ctxProvider";
import useCtx from "./ctxProvider";

interface FolderContentProps {
  content: FolderContentType | undefined;
  isListView: boolean;
  loading: boolean;
}

function FolderContentWithCtx(props: FolderContentProps) {
  const ctx = useCtx();

  const handleContainerClick = (e: Event) => {
    // Clear selection when clicking on empty space
    if (e.target === e.currentTarget) {
      ctx.clearSelection();
    }
  };

  return (
    <div onClick={handleContainerClick} class="min-h-full">
      <Show
        when={!props.loading && props.content}
        fallback={
          <Show when={props.isListView} fallback={<GridSkeleton />}>
            <ListSkeleton />
          </Show>
        }
      >
        <Show
          when={
            props.content!.folderPage.items.length ||
            props.content!.filePage.items.length
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
    </div>
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
