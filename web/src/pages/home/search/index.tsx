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
  const [query, setQuery] = createSignal("");
  const [options, setOptions] = createSignal([]);
  const [loading, setLoading] = createSignal(false);
  const [selectedItem, setSelectedItem] = createSignal(null);

  // Filter and group items based on the search query
  const searchItems = (searchQuery: string) => {
    setLoading(true);

    // Simulate network delay
    setTimeout(() => {
      if (!searchQuery.trim()) {
        setOptions([]);
        setLoading(false);
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
              type: "section",
              name: "Files",
              children: matchingFiles.slice(0, 5).map((file) => {
                const { type: fileType, ...fileRest } = file;
                return {
                  type: "item",
                  ...fileRest,
                  value: file.id,
                  label: file.name,
                };
              }),
            }
          : null;

      const folderSection =
        matchingFolders.length > 0
          ? {
              type: "section",
              name: "Folders",
              children: matchingFolders.slice(0, 3).map((folder) => {
                const { type: folderType, ...folderRest } = folder;
                return {
                  type: "item",
                  ...folderRest,
                  value: folder.id,
                  label: folder.name,
                };
              }),
            }
          : null;

      // Combine sections, filtering out null values
      const result = [fileSection, folderSection].filter(Boolean);

      // Add a single 'See more results...' option if any group exceeds its limit
      const filesExceeded = matchingFiles.length > 5;
      const foldersExceeded = matchingFolders.length > 3;
      if ((filesExceeded || foldersExceeded) && result.length > 0) {
        // Add to the last group
        const lastSection = result[result.length - 1];
        lastSection.children.push({
          type: "item",
          id: "see-more-results",
          name: "See more results...",
          value: "see-more-results",
          label: "See more results...",
          isSeeMore: true,
        });
      }

      setOptions(result);
      setLoading(false);
    }, 700);
  };

  const handleSearchInput = (value: string) => {
    setQuery(value);
    searchItems(value);
  };

  const handleSelectItem = (item) => {
    if (item && item.isSeeMore) {
      // Handle "See more" click (placeholder for now)
      console.log("See more clicked:", item.value);
      return;
    }

    setSelectedItem(item);
    // Here you would typically do something with the selected item
    console.log("Selected item:", item);
  };

  return (
    <KSearch
      options={options()}
      optionValue="value"
      optionLabel="label"
      optionGroupChildren="children"
      optionTextValue="name"
      placeholder="Search..."
      triggerMode="focus"
      debounceOptionsMillisecond={300}
      onInputChange={handleSearchInput}
      onChange={handleSelectItem}
      class=""
      itemComponent={(props) => (
        <KSearch.Item
          item={props.item}
          class="dropdown-item flex items-center gap-2"
        >
          {/* Icon based on item type */}
          <Show when={!props.item.rawValue.isSeeMore}>
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
                  d="M19.5 14.25v-2.625a3.375 3.375 0 0 0-3.375-3.375h-1.5A1.125 1.125 0 0 1 13.5 7.125v-1.5a3.375 3.375 0 0 0-3.375-3.375H8.25m0 12.75h7.5m-7.5 3H12M10.5 2.25H5.625c-.621 0-1.125.504-1.125 1.125v17.25c0 .621.504 1.125 1.125 1.125h12.75c.621 0 1.125-.504 1.125-1.125V11.25a9 9 0 0 0-9-9Z"
                />
              </svg>
            </Show>
          </Show>
          <Show when={props.item.rawValue.isSeeMore}>
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
                d="M8.25 4.5l7.5 7.5-7.5 7.5"
              />
            </svg>
          </Show>

          {/* Item label with different styling for "See more" */}
          <KSearch.ItemLabel
            class={
              props.item.rawValue.isSeeMore
                ? "text-primary text-sm font-medium"
                : ""
            }
          >
            {props.item.rawValue.label}
          </KSearch.ItemLabel>
        </KSearch.Item>
      )}
      sectionComponent={(props) => {
        // In production, you'd use a better approach, but this works for now
        const name = props.item
          ? props.item.children?.[0]?.type === "folder"
            ? "Folders"
            : "Files"
          : "Section";

        return (
          <KSearch.Section class="px-3 py-1 text-xs text-neutral-light font-semibold uppercase">
            {name}
          </KSearch.Section>
        );
      }}
    >
      <KSearch.Control
        aria-label="Search"
        class="flex items-center bg-white border rounded-md input input-b-std"
      >
        <KSearch.Indicator
          loadingComponent={
            <KSearch.Icon class="flex-center pl-2">
              <svg
                class="size-4 text-neutral-light animate-spin"
                xmlns="http://www.w3.org/2000/svg"
                fill="none"
                viewBox="0 0 24 24"
              >
                <circle
                  class="opacity-25"
                  cx="12"
                  cy="12"
                  r="10"
                  stroke="currentColor"
                  stroke-width="4"
                ></circle>
                <path
                  class="opacity-75"
                  fill="currentColor"
                  d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                ></path>
              </svg>
            </KSearch.Icon>
          }
        >
          <KSearch.Icon class="flex-center pl-2">
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
        <KSearch.Input class="w-full border-none focus:ring-0 outline-none py-2" />
      </KSearch.Control>

      <KSearch.Portal>
        <KSearch.Content class="bg-white rounded-md shadow-md border border-border-strong mt-1 py-1 max-h-[300px] overflow-y-auto z-50">
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
