import { createFolder } from "@sv/apis/media";
import { COMMON_ERR_KEYS } from "@sv/utils/errors";
import { Accessor, createEffect, createSignal, Setter } from "solid-js";

interface Props {
  isModalOpen: Accessor<boolean>;
  setIsModalOpen: Setter<boolean>;
}

function useVM(props: Props) {
  const [folderName, setFolderName] = createSignal("");
  const [isCreating, setIsCreating] = createSignal(false);
  const [error, setError] = createSignal("");

  // Reset the folder name when the modal is closed
  createEffect(() => {
    if (!props.isModalOpen()) {
      setFolderName("");
      setError("");
    }
  });

  const handleFolderNameChange = (name: string) => {
    setFolderName(name);
    isInvalidFolderName();
  };

  const isInvalidFolderName = (): boolean => {
    if (!folderName() || folderName().trim() === "") {
      setError("Folder name is required");
      return true;
    }

    setError("");
    return false;
  };

  const handleCreateFolder = async (parentFolderId?: string) => {
    if (isInvalidFolderName()) {
      return;
    }

    setIsCreating(true);

    try {
      await createFolder(parentFolderId, folderName().trim());
      setFolderName(""); // Reset the name after successful creation
    } catch (err) {
      if (err instanceof Error && err.message === COMMON_ERR_KEYS.DUPLICATE) {
        setError("Folder already exists");
      } else {
        setError("Failed to create folder");
      }
    } finally {
      setIsCreating(false);
    }

    if (!error()) {
      props.setIsModalOpen(false);
    }
  };

  return {
    folderName,
    isCreating,
    error,
    handleFolderNameChange,
    handleCreateFolder,
  };
}

export default useVM;
