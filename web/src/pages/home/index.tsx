import { createResource, Show, For, createSignal } from "solid-js";
import StorageMeter from "@sv/components/StorageMeter";
import ItemContextMenu from "@sv/components/ItemContextMenu";
import EmptyState from "@sv/components/EmptyState";
import SkeletonGrid from "@sv/components/SkeletonGrid";
import FileCard from "@sv/components/file/FileCard";
import FileListItem from "@sv/components/file/FileListItem";
import FolderCard from "@sv/components/folder/FolderCard";
import FolderListItem from "@sv/components/folder/FolderListItem";
import { fetchRootContent } from "@sv/apis/media";
import type { FolderContent } from "@sv/apis/media/models";
import { Button } from "@kobalte/core/button";

// Dummy storage usage
const storageUsed = 2.5; // GB
const storageTotal = 10; // GB

export default function Home() {
  const [rootContent] = createResource<FolderContent>(fetchRootContent);
  const [viewMode, setViewMode] = createSignal<"grid" | "list">("grid");

  return (
    <>
      {/* Page title and breadcrumbs */}
      <div class="flex items-center justify-between">
        <div>
          <h1 class="text-2xl font-semibold text-gray-800">Home</h1>
          <div class="text-sm text-gray-500 mt-1">
            All your files in one secure place
          </div>
        </div>
        <div class="flex items-center gap-2">
          <button
            class={`p-2 rounded relative group`}
            title="Grid View"
            aria-label="Grid View"
            style={{
              background:
                viewMode() === "grid"
                  ? "linear-gradient(90deg,#6366f1 0%,#38bdf8 100%)"
                  : undefined,
              color: viewMode() === "grid" ? "#2563eb" : undefined,
            }}
            onClick={() => setViewMode("grid")}
          >
            <span class="material-symbols-outlined">grid_view</span>
          </button>
          <button
            class={`p-2 rounded relative group`}
            title="List View"
            aria-label="List View"
            style={{
              background:
                viewMode() === "list"
                  ? "linear-gradient(90deg,#6366f1 0%,#38bdf8 100%)"
                  : undefined,
              color: viewMode() === "list" ? "#2563eb" : undefined,
            }}
            onClick={() => setViewMode("list")}
          >
            <span class="material-symbols-outlined">view_list</span>
          </button>
        </div>
      </div>

      {/* Storage meter and actions */}
      <div class="bg-white rounded-lg shadow-sm border border-gray-200 p-4">
        <div class="flex flex-col md:flex-row md:items-center gap-4">
          <StorageMeter used={storageUsed} total={storageTotal} />
          <div class="flex flex-wrap gap-3">
            <Button class="btn btn-primary flex items-center gap-2">
              <span class="material-symbols-outlined text-sm">upload</span>
              <span>Upload</span>
            </Button>
            <Button class="btn btn-outline flex items-center gap-2">
              <span class="material-symbols-outlined text-sm">
                create_new_folder
              </span>
              <span>New Folder</span>
            </Button>
          </div>
        </div>
      </div>

      {/* Files & Folders */}
      <Show
        when={rootContent.loading}
        fallback={
          <Show
            when={
              (rootContent() &&
                (rootContent()?.folderPage?.items?.length ?? 0) > 0) ||
              (rootContent()?.filePage?.items?.length ?? 0) > 0
            }
            fallback={<EmptyState />}
          >
            <Show
              when={viewMode() === "grid"}
              fallback={
                <div class="bg-white rounded-lg border border-gray-200 shadow-sm divide-y">
                  {/* Folder list items */}
                  <For each={rootContent()?.folderPage?.items}>
                    {(folder) => (
                      <FolderListItem name={folder.name}>
                        <ItemContextMenu
                          type="folder"
                          id={folder.id}
                          name={folder.name}
                        />
                      </FolderListItem>
                    )}
                  </For>

                  {/* File list items */}
                  <For each={rootContent()?.filePage?.items}>
                    {(file) => (
                      <FileListItem
                        name={file.name}
                        size={file.size}
                        updatedAt={file.updatedAt}
                      >
                        <ItemContextMenu
                          type="file"
                          id={file.id}
                          name={file.name}
                        />
                      </FileListItem>
                    )}
                  </For>
                </div>
              }
            >
              <div class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 gap-4">
                {/* Folders */}
                <For each={rootContent()?.folderPage?.items}>
                  {(folder) => (
                    <FolderCard name={folder.name}>
                      <ItemContextMenu
                        type="folder"
                        id={folder.id}
                        name={folder.name}
                      />
                    </FolderCard>
                  )}
                </For>

                {/* Files */}
                <For each={rootContent()?.filePage?.items}>
                  {(file) => (
                    <FileCard
                      name={file.name}
                      size={file.size}
                      updatedAt={file.updatedAt}
                    >
                      <ItemContextMenu
                        type="file"
                        id={file.id}
                        name={file.name}
                      />
                    </FileCard>
                  )}
                </For>
              </div>
            </Show>
          </Show>
        }
      >
        {/* Skeleton loader */}
        <SkeletonGrid />
      </Show>
    </>
  );
}
