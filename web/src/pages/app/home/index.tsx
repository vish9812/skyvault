import { Show, For } from "solid-js";
import { Button } from "@kobalte/core/button";
import { A } from "@solidjs/router";
import FolderContent from "@sv/components/folder-content";
import useViewModel from "./useViewModel";

const folderParentPath = ["Home", "Folder 1", "Folder 11"];
const currentFolder = "Folder 111";

export default function Home() {
  const { isListView, handleListViewChange, folderContentRes } = useViewModel();

  return (
    <>
      <div class="flex items-center justify-between">
        {/* Breadcrumbs */}
        <div class="text-primary">
          <For each={folderParentPath}>
            {(folder) => (
              <span>
                <A href="/" class="link">
                  {folder}
                </A>
                {" / "}
              </span>
            )}
          </For>
          <span class="font-bold text-neutral">{currentFolder}</span>
        </div>
        <div>
          <Button class="btn btn-ghost" onClick={handleListViewChange}>
            <Show
              when={isListView()}
              fallback={
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke-width="1.5"
                  stroke="currentColor"
                  class="size-6"
                >
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    d="M3.75 6A2.25 2.25 0 0 1 6 3.75h2.25A2.25 2.25 0 0 1 10.5 6v2.25a2.25 2.25 0 0 1-2.25 2.25H6a2.25 2.25 0 0 1-2.25-2.25V6ZM3.75 15.75A2.25 2.25 0 0 1 6 13.5h2.25a2.25 2.25 0 0 1 2.25 2.25V18a2.25 2.25 0 0 1-2.25 2.25H6A2.25 2.25 0 0 1 3.75 18v-2.25ZM13.5 6a2.25 2.25 0 0 1 2.25-2.25H18A2.25 2.25 0 0 1 20.25 6v2.25A2.25 2.25 0 0 1 18 10.5h-2.25a2.25 2.25 0 0 1-2.25-2.25V6ZM13.5 15.75a2.25 2.25 0 0 1 2.25-2.25H18a2.25 2.25 0 0 1 2.25 2.25V18A2.25 2.25 0 0 1 18 20.25h-2.25A2.25 2.25 0 0 1 13.5 18v-2.25Z"
                  />
                </svg>
              }
            >
              <svg
                xmlns="http://www.w3.org/2000/svg"
                fill="none"
                viewBox="0 0 24 24"
                stroke-width="1.5"
                stroke="currentColor"
                class="size-6"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  d="M8.25 6.75h12M8.25 12h12m-12 5.25h12M3.75 6.75h.007v.008H3.75V6.75Zm.375 0a.375.375 0 1 1-.75 0 .375.375 0 0 1 .75 0ZM3.75 12h.007v.008H3.75V12Zm.375 0a.375.375 0 1 1-.75 0 .375.375 0 0 1 .75 0Zm-.375 5.25h.007v.008H3.75v-.008Zm.375 0a.375.375 0 1 1-.75 0 .375.375 0 0 1 .75 0Z"
                />
              </svg>
            </Show>
          </Button>
        </div>
      </div>

      {/* Files & Folders */}
      <FolderContent
        loading={folderContentRes.loading}
        content={folderContentRes.latest}
        isListView={isListView()}
      />
    </>
  );
}
