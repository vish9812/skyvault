import { createEffect, createSignal, ParentProps, useContext } from "solid-js";
import CTX from "./ctx";
import { createFolder } from "@sv/apis/media";
import { VALIDATIONS } from "@sv/utils/consts";

export function CtxProvider(props: ParentProps) {
  const [isCreateFolderModalOpen, setIsCreateFolderModalOpen] =
    createSignal(false);
  const [createFolderName, setCreateFolderName] = createSignal("");
  const [isCreating, setIsCreating] = createSignal(false);
  const [error, setError] = createSignal("");

  // Reset the folder name when the modal is closed
  createEffect(() => {
    if (!isCreateFolderModalOpen()) {
      setCreateFolderName("");
      setError("");
    }
  });

  const handleCreateFolder = async (parentFolderId?: string) => {
    if (!createFolderName() || createFolderName().trim() === "") {
      setError("Folder name is required");
      return;
    }

    if (createFolderName().length > VALIDATIONS.MAX_LENGTH) {
      setError("Folder name is too long");
      return;
    }

    setIsCreating(true);
    setError("");

    try {
      await createFolder(parentFolderId, createFolderName());
      setCreateFolderName(""); // Reset the name after successful creation
    } catch (err) {
      setError("Failed to create folder");
    } finally {
      setIsCreating(false);
    }

    if (!error()) {
      setIsCreateFolderModalOpen(false);
    }
  };

  const val = {
    isFolderModalOpen: isCreateFolderModalOpen,
    setIsFolderModalOpen: setIsCreateFolderModalOpen,
    folderName: createFolderName,
    setFolderName: setCreateFolderName,
    isCreating,
    error,
    handleCreateFolder,
  };

  return <CTX.Provider value={val}>{props.children}</CTX.Provider>;
}

function useCtx() {
  const ctx = useContext(CTX);
  if (!ctx) {
    throw new Error("createUpload: useCtx must be used within a CtxProvider");
  }
  return ctx;
}

export default useCtx;
