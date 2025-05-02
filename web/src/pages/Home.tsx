import { createResource, Show, For, createSignal } from "solid-js";
import { Meter } from "@kobalte/core/meter";
import { Button } from "@kobalte/core/button";
import { DropdownMenu } from "@kobalte/core/dropdown-menu";
import AppLayout from "@sv/components/AppLayout";

// Dummy storage usage
const storageUsed = 2.5; // GB
const storageTotal = 10; // GB

// Type for API response
interface FolderContent {
  folders: { id: number; name: string }[];
  files: { id: number; name: string }[];
}

async function fetchRootContent(): Promise<FolderContent> {
  const token = localStorage.getItem("token");
  const res = await fetch(
    "http://localhost:8090/api/v1/media/folders/0/content",
    {
      headers: { Authorization: `Bearer ${token}` },
    }
  );
  if (!res.ok) throw new Error("Failed to fetch folder content");
  return res.json();
}

// Context menu for a file or folder item
const ItemContextMenu = (props: {
  type: "file" | "folder";
  id: number;
  name: string;
}) => {
  return (
    <DropdownMenu>
      <DropdownMenu.Trigger class="text-gray-400 hover:text-gray-600 rounded-full p-1 absolute top-2 right-2 opacity-0 group-hover:opacity-100 transition-opacity">
        <span class="material-symbols-outlined text-base">more_vert</span>
      </DropdownMenu.Trigger>
      <DropdownMenu.Content class="bg-white rounded-lg shadow-lg border border-gray-200 py-1 min-w-[180px]">
        <DropdownMenu.Item class="px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 cursor-pointer flex items-center gap-2">
          <span class="material-symbols-outlined text-gray-500 text-sm">
            drive_file_rename_outline
          </span>
          <span>Rename</span>
        </DropdownMenu.Item>
        <DropdownMenu.Item class="px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 cursor-pointer flex items-center gap-2">
          <span class="material-symbols-outlined text-gray-500 text-sm">
            drive_file_move
          </span>
          <span>Move</span>
        </DropdownMenu.Item>
        <DropdownMenu.Item class="px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 cursor-pointer flex items-center gap-2">
          <span class="material-symbols-outlined text-gray-500 text-sm">
            share
          </span>
          <span>Share</span>
        </DropdownMenu.Item>
        <DropdownMenu.Item class="px-4 py-2 text-sm text-red-600 hover:bg-gray-100 cursor-pointer flex items-center gap-2">
          <span class="material-symbols-outlined text-red-500 text-sm">
            delete
          </span>
          <span>Delete</span>
        </DropdownMenu.Item>
      </DropdownMenu.Content>
    </DropdownMenu>
  );
};

export default function Home() {
  const [rootContent] = createResource(fetchRootContent);
  const [viewMode, setViewMode] = createSignal<"grid" | "list">("grid");

  return (
    <AppLayout>
      <div class="max-w-6xl mx-auto space-y-6">
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
              class={`p-2 rounded ${
                viewMode() === "grid"
                  ? "bg-primary/20 text-primary"
                  : "text-gray-500 hover:bg-primary/10"
              }`}
              onClick={() => setViewMode("grid")}
            >
              <span class="material-symbols-outlined">grid_view</span>
            </button>
            <button
              class={`p-2 rounded ${
                viewMode() === "list"
                  ? "bg-primary/20 text-primary"
                  : "text-gray-500 hover:bg-primary/10"
              }`}
              onClick={() => setViewMode("list")}
            >
              <span class="material-symbols-outlined">view_list</span>
            </button>
          </div>
        </div>

        {/* Storage meter and actions */}
        <div class="bg-white rounded-lg shadow-sm border border-gray-200 p-4">
          <div class="flex flex-col md:flex-row md:items-center gap-4">
            <div class="flex-1">
              <div class="flex justify-between mb-1">
                <span class="text-sm font-medium text-gray-700">
                  Storage Usage
                </span>
                <span class="text-sm text-gray-500">
                  {storageUsed} GB of {storageTotal} GB used
                </span>
              </div>
              <Meter
                value={storageUsed}
                minValue={0}
                maxValue={storageTotal}
                class="w-full"
              >
                <Meter.Track class="h-2 bg-gray-200 rounded-full overflow-hidden">
                  <Meter.Fill class="bg-primary h-full rounded-full transition-all" />
                </Meter.Track>
              </Meter>
            </div>
            <div class="flex flex-wrap gap-2">
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
                rootContent() &&
                (rootContent()!.folders.length > 0 ||
                  rootContent()!.files.length > 0)
              }
              fallback={
                <div class="text-center py-16 bg-white rounded-lg border border-gray-200 shadow-sm">
                  <div class="inline-flex justify-center items-center w-16 h-16 rounded-full bg-primary/20 text-primary mb-4">
                    <span class="material-symbols-outlined text-3xl">
                      cloud_upload
                    </span>
                  </div>
                  <h3 class="text-lg font-medium text-gray-900">
                    No files or folders yet
                  </h3>
                  <p class="mt-2 text-sm text-gray-500 max-w-md mx-auto">
                    Upload files or create folders to get started with your
                    secure cloud storage.
                  </p>
                  <div class="mt-6">
                    <Button class="btn btn-primary">Upload Files</Button>
                  </div>
                </div>
              }
            >
              <Show
                when={viewMode() === "grid"}
                fallback={
                  <div class="bg-white rounded-lg border border-gray-200 shadow-sm divide-y">
                    {/* Folder list items */}
                    <For each={rootContent()?.folders}>
                      {(folder) => (
                        <div class="flex items-center py-3 px-4 hover:bg-gray-50 group relative">
                          <span class="material-symbols-outlined text-primary mr-3">
                            folder
                          </span>
                          <div class="flex-1 min-w-0">
                            <div class="text-sm font-medium text-gray-900 truncate">
                              {folder.name}
                            </div>
                          </div>
                          <ItemContextMenu
                            type="folder"
                            id={folder.id}
                            name={folder.name}
                          />
                        </div>
                      )}
                    </For>

                    {/* File list items */}
                    <For each={rootContent()?.files}>
                      {(file) => (
                        <div class="flex items-center py-3 px-4 hover:bg-gray-50 group relative">
                          <span class="material-symbols-outlined text-gray-400 mr-3">
                            description
                          </span>
                          <div class="flex-1 min-w-0">
                            <div class="text-sm font-medium text-gray-900 truncate">
                              {file.name}
                            </div>
                          </div>
                          <ItemContextMenu
                            type="file"
                            id={file.id}
                            name={file.name}
                          />
                        </div>
                      )}
                    </For>
                  </div>
                }
              >
                <div class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 gap-4">
                  {/* Folders */}
                  <For each={rootContent()?.folders}>
                    {(folder) => (
                      <div class="bg-white rounded-lg shadow-sm border border-gray-200 p-4 hover:shadow-md transition-shadow relative group">
                        <div class="flex flex-col items-center">
                          <span class="material-symbols-outlined text-4xl text-primary mb-2">
                            folder
                          </span>
                          <div class="font-medium text-sm text-center text-gray-900 truncate w-full">
                            {folder.name}
                          </div>
                        </div>
                        <ItemContextMenu
                          type="folder"
                          id={folder.id}
                          name={folder.name}
                        />
                      </div>
                    )}
                  </For>

                  {/* Files */}
                  <For each={rootContent()?.files}>
                    {(file) => (
                      <div class="bg-white rounded-lg shadow-sm border border-gray-200 p-4 hover:shadow-md transition-shadow relative group">
                        <div class="flex flex-col items-center">
                          <span class="material-symbols-outlined text-4xl text-gray-400 mb-2">
                            description
                          </span>
                          <div class="font-medium text-sm text-center text-gray-900 truncate w-full">
                            {file.name}
                          </div>
                        </div>
                        <ItemContextMenu
                          type="file"
                          id={file.id}
                          name={file.name}
                        />
                      </div>
                    )}
                  </For>
                </div>
              </Show>
            </Show>
          }
        >
          {/* Skeleton loader */}
          <div class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 gap-4">
            <For each={Array(10)}>
              {() => (
                <div class="bg-white rounded-lg shadow-sm border border-gray-200 p-4">
                  <div class="flex flex-col items-center">
                    <div class="w-10 h-10 bg-gray-200 rounded mb-3 animate-pulse"></div>
                    <div class="h-4 bg-gray-200 rounded w-20 animate-pulse"></div>
                  </div>
                </div>
              )}
            </For>
          </div>
        </Show>
      </div>
    </AppLayout>
  );
}
