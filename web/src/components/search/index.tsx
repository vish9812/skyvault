import { Button } from "@kobalte/core/button";
import { Search as KSearch } from "@kobalte/core/search";
import { createSignal, createEffect, Show, For } from "solid-js";

// Mock data for demonstration
const MOCK_FILES = [
  { id: "f1", name: "Project Proposal.pdf", type: "file" },
  { id: "f2", name: "Budget 2023.xlsx", type: "file" },
  { id: "f3", name: "Meeting Notes.docx", type: "file" },
  { id: "f4", name: "Presentation.pptx", type: "file" },
  { id: "f5", name: "Logo.png", type: "file" },
  { id: "f6", name: "Product Roadmap.pdf", type: "file" },
  { id: "f7", name: "User Research.docx", type: "file" },
  { id: "f8", name: "Marketing Plan.pdf", type: "file" },
  { id: "f9", name: "Personal File.txt", type: "file" },
];

const MOCK_FOLDERS = [
  { id: "d1", name: "Documents", type: "folder" },
  { id: "d2", name: "Images", type: "folder" },
  { id: "d3", name: "Projects", type: "folder" },
  { id: "d4", name: "Archive", type: "folder" },
  { id: "d5", name: "Personal", type: "folder" },
  { id: "d6", name: "Work Templates", type: "folder" },
  { id: "d7", name: "Personal Templates", type: "folder" },
];

// Combined data
const ALL_ITEMS = [...MOCK_FILES, ...MOCK_FOLDERS];

function Search() {
  let inputRef: HTMLInputElement | undefined;
  // const [query, setQuery] = createSignal({
  //   id: "",
  //   name: "",
  //   type: "",
  // });
  const [options, setOptions] = createSignal([]);
  // const [selectedItem, setSelectedItem] = createSignal(null);

  const handleSearchClose = () => {
    if (inputRef) {
      inputRef.value = "";
    }
    setOptions([]);

    // setSelectedItem(null);
  };

  // Filter and group items based on the search query
  const searchItems = (searchQuery: string) => {
    // Simulate network delay
    setTimeout(() => {
      if (!searchQuery.trim()) {
        setOptions([]);
        return;
      }

      const lowercaseQuery = searchQuery.toLowerCase();

      // Filter files and folders that match the query
      const matchingFiles = MOCK_FILES.filter((file) =>
        file.name.toLowerCase().includes(lowercaseQuery)
      );

      const matchingFolders = MOCK_FOLDERS.filter((folder) =>
        folder.name.toLowerCase().includes(lowercaseQuery)
      );

      // Create sections for files and folders
      const fileSection =
        matchingFiles.length > 0
          ? {
              name: "Files",
              children: matchingFiles.slice(0, 5),
            }
          : null;

      const folderSection =
        matchingFolders.length > 0
          ? {
              name: "Folders",
              children: matchingFolders.slice(0, 3),
            }
          : null;

      // Combine sections, filtering out null values
      const result = [fileSection, folderSection].filter(Boolean);

      // Add a single 'See more results...' option if any group exceeds its limit
      const filesExceeded = matchingFiles.length > 5;
      const foldersExceeded = matchingFolders.length > 3;
      if ((filesExceeded || foldersExceeded) && result.length > 0) {
        result.push({
          name: "",
          children: [
            {
              type: "see-more",
              id: "see-more-results",
              name: "See more results...",
            },
          ],
        });
      }

      setOptions(result);
    }, 300);
  };

  const handleSearchInput = (value) => {
    // setQuery(value);
    searchItems(value);
  };

  // const handleSelectItem = (item) => {
  //   if (item && item.isSeeMore) {
  //     // Handle "See more" click (placeholder for now)
  //
  //     return;
  //   }

  //   setSelectedItem(item);
  //   // Here you would typically do something with the selected item
  //
  // };

  return (
    <KSearch
      // value={query()}
      options={options()}
      optionLabel="name"
      optionGroupChildren="children"
      optionTextValue="name"
      placeholder="Search"
      triggerMode="focus"
      debounceOptionsMillisecond={500}
      closeOnSelection={true}
      onInputChange={handleSearchInput}
      // onChange={handleSelectItem}
      itemComponent={(props) => (
        <KSearch.Item
          item={props.item}
          class="flex items-center justify-between"
        >
          {/* Icon based on item type */}
          <div class="flex items-center gap-2 dropdown-item">
            <Show when={props.item.rawValue.type !== "see-more"}>
              <Show when={props.item.rawValue.type === "folder"}>
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke-width="1.5"
                  stroke="currentColor"
                  class="size-4 text-neutral-light"
                >
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    d="M2.25 12.75V12A2.25 2.25 0 0 1 4.5 9.75h15A2.25 2.25 0 0 1 21.75 12v.75m-8.69-6.44-2.12-2.12a1.5 1.5 0 0 0-1.061-.44H4.5A2.25 2.25 0 0 0 2.25 6v12a2.25 2.25 0 0 0 2.25 2.25h15A2.25 2.25 0 0 0 21.75 18V9a2.25 2.25 0 0 0-2.25-2.25h-5.379a1.5 1.5 0 0 1-1.06-.44Z"
                  />
                </svg>
              </Show>
              <Show when={props.item.rawValue.type === "file"}>
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke-width="1.5"
                  stroke="currentColor"
                  class="size-4 text-neutral-light"
                >
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    d="M19.5 14.25v-2.625a3.375 3.375 0 0 0-3.375-3.375h-1.5A1.125 1.125 0 0 1 13.5 7.125v-1.5a3.375 3.375 0 0 0-3.375-3.375H8.25m2.25 0H5.625c-.621 0-1.125.504-1.125 1.125v17.25c0 .621.504 1.125 1.125 1.125h12.75c.621 0 1.125-.504 1.125-1.125V11.25a9 9 0 0 0-9-9Z"
                  />
                </svg>
              </Show>
            </Show>
            <Show when={props.item.rawValue.type === "see-more"}>
              <svg
                xmlns="http://www.w3.org/2000/svg"
                fill="none"
                viewBox="0 0 24 24"
                stroke-width="1.5"
                stroke="currentColor"
                class="size-4 text-neutral-light"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  d="m8.25 4.5 7.5 7.5-7.5 7.5"
                />
              </svg>
            </Show>

            {/* Item label with different styling for "See more" */}
            <KSearch.ItemLabel
              class={
                props.item.rawValue.type === "see-more"
                  ? "link text-sm font-medium"
                  : ""
              }
            >
              {props.item.rawValue.name}
            </KSearch.ItemLabel>
          </div>
          {/* Download icon */}
          <button
            type="button"
            class="p-1 rounded-full hover:bg-secondary-lighter"
          >
            <svg
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 24 24"
              stroke-width="1.5"
              stroke="currentColor"
              class="size-4 text-primary"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                d="M3 16.5v2.25A2.25 2.25 0 0 0 5.25 21h13.5A2.25 2.25 0 0 0 21 18.75V16.5M16.5 12 12 16.5m0 0L7.5 12m4.5 4.5V3"
              />
            </svg>
          </button>
        </KSearch.Item>
      )}
      sectionComponent={(props) => (
        <KSearch.Section class="px-3 py-1 text-xs text-neutral-light font-semibold uppercase">
          {props.section.rawValue.name}
        </KSearch.Section>
      )}
    >
      <KSearch.Control
        aria-label="Search"
        class="flex-center w-[200px] md:w-[350px] mx-4 bg-white input-b-std border rounded-md focus:outline-none"
      >
        <KSearch.Indicator
          loadingComponent={
            <KSearch.Icon class="flex-center p-1">
              <svg
                xmlns="http://www.w3.org/2000/svg"
                fill="none"
                viewBox="0 0 24 24"
                stroke-width="1.5"
                stroke="currentColor"
                class="size-4 text-neutral-light animate-spin"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  d="M16.023 9.348h4.992v-.001M2.985 19.644v-4.992m0 0h4.992m-4.993 0 3.181 3.183a8.25 8.25 0 0 0 13.803-3.7M4.031 9.865a8.25 8.25 0 0 1 13.803-3.7l3.181 3.182m0-4.991v4.99"
                />
              </svg>
            </KSearch.Icon>
          }
        >
          <KSearch.Icon class="flex-center p-1">
            <svg
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 24 24"
              stroke-width="1.5"
              stroke="currentColor"
              class="size-4 text-neutral-light"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                d="m21 21-5.197-5.197m0 0A7.5 7.5 0 1 0 5.196 5.196a7.5 7.5 0 0 0 10.607 10.607Z"
              />
            </svg>
          </KSearch.Icon>
        </KSearch.Indicator>
        <KSearch.Input ref={inputRef} class="w-full outline-none py-2" />
        {/* Close button */}
        <Button
          class="p-1"
          onPointerDown={(e) => e.stopPropagation()}
          onClick={handleSearchClose}
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
            stroke-width="1.5"
            stroke="currentColor"
            class="size-5 text-neutral-light"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              d="M6 18L18 6M6 6l12 12"
            />
          </svg>
        </Button>
      </KSearch.Control>

      <KSearch.Portal>
        <KSearch.Content
          class="bg-white rounded-md border border-border-strong py-1 max-h-[500px] overflow-y-auto z-20"
          onCloseAutoFocus={(e) => e.preventDefault()}
        >
          <KSearch.Listbox />
          <KSearch.NoResult class="px-3 py-2 text-sm text-neutral-light">
            No results found
          </KSearch.NoResult>
        </KSearch.Content>
      </KSearch.Portal>
    </KSearch>
  );
}

export default Search;
