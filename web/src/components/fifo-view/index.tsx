// FIFO is a short term for Files and Folders

import { createEffect, Show } from "solid-js";
import EmptyState from "./empty-state";
import { FolderContent } from "@sv/apis/media/models";
import GridSkeleton from "./grid-skeleton";
import ListSkeleton from "./list-skeleton";
import ListView from "./list-view";
import GridView from "./grid-view";

interface FifoViewProps {
  fifo: FolderContent | undefined;
  isListView: boolean;
  loading: boolean;
}

function FifoView(props: FifoViewProps) {
  createEffect(() => {
    console.log("loading", props.loading);
    console.log("isListView", props.isListView);
  });

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
          (props.fifo?.folderPage?.items?.length ?? 0) > 0 ||
          (props.fifo?.filePage?.items?.length ?? 0) > 0
        }
        fallback={<EmptyState />}
      >
        <Show
          when={props.isListView}
          fallback={<GridView content={props.fifo!} />}
        >
          <ListView content={props.fifo!} />
        </Show>
      </Show>
    </Show>
  );
}

export default FifoView;
