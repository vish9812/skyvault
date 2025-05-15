// FIFO is a short term for Files and Folders

import { createEffect, For, Show } from "solid-js";
import EmptyState from "./empty-state";
import { FolderContent } from "@sv/apis/media/models";
import Row from "./row";
import Card from "./card";
import CardsSkeleton from "./cards-skeleton";
import RowsSkeleton from "./rows-skeleton";

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
        <Show when={props.isListView} fallback={<CardsSkeleton />}>
          <RowsSkeleton />
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
          fallback={
            <div class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 gap-4">
              {/* Folders */}
              <For each={props.fifo?.folderPage?.items}>
                {(folder) => <Card type="folder" name={folder.name} />}
              </For>

              {/* Files */}
              <For each={props.fifo?.filePage?.items}>
                {(file) => <Card type="file" name={file.name} />}
              </For>
            </div>
          }
        >
          <div>
            {/* Folder list items */}
            <For each={props.fifo?.folderPage?.items}>
              {(folder) => <Row type="folder" name={folder.name} />}
            </For>

            {/* File list items */}
            <For each={props.fifo?.filePage?.items}>
              {(file) => <Row type="file" name={file.name} />}
            </For>
          </div>
        </Show>
      </Show>
    </Show>
  );
}

export default FifoView;
