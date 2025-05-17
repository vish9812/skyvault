import { DropdownMenu } from "@kobalte/core/dropdown-menu";
import useViewModel from "./useViewModel";
import { CreateFolderModal } from "./createFolderModal";

function CreateUpload() {
  const vm = useViewModel();

  const handleCreateFolderOption = () => {
    vm.setIsFolderModalOpen(true);
  };

  const handleCreateFolder = () => {
    vm.handleCreateFolder();
    if (!vm.error()) {
      vm.setIsFolderModalOpen(false);
    }
  };

  return (
    <>
      <DropdownMenu>
        <DropdownMenu.Trigger class="p-2 btn btn-gradient btn-gradient-d-expanded max-md:rounded-full md:w-full md:mb-4 md:mt-2">
          {/* Desktop trigger */}
          <span class="hidden md:block">
            <span class="flex-center gap-2">
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
                  d="M12 9v6m3-3H9m12 0a9 9 0 1 1-18 0 9 9 0 0 1 18 0Z"
                />
              </svg>
              Create or Upload
            </span>
          </span>

          {/* Mobile trigger */}
          <div class="md:hidden">
            <svg
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 24 24"
              stroke-width="1.5"
              stroke="currentColor"
              class="size-6 md:hidden"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                d="M12 4.5v15m7.5-7.5h-15"
              />
            </svg>
          </div>
        </DropdownMenu.Trigger>
        <DropdownMenu.Portal>
          <DropdownMenu.Content class="bg-white rounded-lg shadow-md border border-border-strong min-w-[180px] py-2">
            <DropdownMenu.Item
              class="dropdown-item"
              onSelect={handleCreateFolderOption}
            >
              <div class="flex items-center gap-2">
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
                    d="M2.25 12.75V12A2.25 2.25 0 0 1 4.5 9.75h15A2.25 2.25 0 0 1 21.75 12v.75m-8.69-6.44-2.12-2.12a1.5 1.5 0 0 0-1.061-.44H4.5A2.25 2.25 0 0 0 2.25 6v12a2.25 2.25 0 0 0 2.25 2.25h15A2.25 2.25 0 0 0 21.75 18V9a2.25 2.25 0 0 0-2.25-2.25h-5.379a1.5 1.5 0 0 1-1.06-.44Z"
                  />
                </svg>
                <span>Create Folder</span>
              </div>
            </DropdownMenu.Item>
            <DropdownMenu.Separator class="border-border-strong my-2" />
            <DropdownMenu.Item class="dropdown-item">
              <div class="flex items-center gap-2">
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
                    d="M19.5 14.25v-2.625a3.375 3.375 0 0 0-3.375-3.375h-1.5A1.125 1.125 0 0 1 13.5 7.125v-1.5a3.375 3.375 0 0 0-3.375-3.375H8.25m3.75 9v6m3-3H9m1.5-12H5.625c-.621 0-1.125.504-1.125 1.125v17.25c0 .621.504 1.125 1.125 1.125h12.75c.621 0 1.125-.504 1.125-1.125V11.25a9 9 0 0 0-9-9Z"
                  />
                </svg>

                <span>Upload Files</span>
              </div>
            </DropdownMenu.Item>
            <DropdownMenu.Item class="dropdown-item">
              <div class="flex items-center gap-2">
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
                    d="M12 10.5v6m3-3H9m4.06-7.19-2.12-2.12a1.5 1.5 0 0 0-1.061-.44H4.5A2.25 2.25 0 0 0 2.25 6v12a2.25 2.25 0 0 0 2.25 2.25h15A2.25 2.25 0 0 0 21.75 18V9a2.25 2.25 0 0 0-2.25-2.25h-5.379a1.5 1.5 0 0 1-1.06-.44Z"
                  />
                </svg>
                <span>Upload Folder</span>
              </div>
            </DropdownMenu.Item>
            <DropdownMenu.Arrow />
          </DropdownMenu.Content>
        </DropdownMenu.Portal>
      </DropdownMenu>

      <CreateFolderModal
        isOpen={vm.isFolderModalOpen()}
        onClose={() => vm.setIsFolderModalOpen(false)}
        folderName={vm.folderName()}
        setFolderName={vm.setFolderName}
        onCreateFolder={handleCreateFolder}
        isCreating={vm.isCreating()}
        error={vm.error()}
      />
    </>
  );
}

export default CreateUpload;
