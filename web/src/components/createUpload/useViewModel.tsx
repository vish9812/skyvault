import { createSignal } from "solid-js";
import { createFolder } from "@sv/apis/media";

function useViewModel() {
  const [isFolderModalOpen, setIsFolderModalOpen] = createSignal(false);
  const [folderName, setFolderName] = createSignal("");
  const [isCreating, setIsCreating] = createSignal(false);
  const [error, setError] = createSignal<string | null>(null);

  const handleCreateFolder = async (parentFolderId?: string) => {
    if (!folderName()) {
      setError("Folder name is required");
      return;
    }

    setIsCreating(true);
    setError(null);

    try {
      await createFolder(parentFolderId, folderName());
      setFolderName(""); // Reset the name after successful creation
    } catch (err) {
      setError("Failed to create folder");
      console.error(err);
    } finally {
      setIsCreating(false);
    }
  };

  return {
    isFolderModalOpen,
    setIsFolderModalOpen,
    folderName,
    setFolderName,
    isCreating,
    error,
    handleCreateFolder,
  };
}

export default useViewModel;
