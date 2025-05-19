import { createEffect, createSignal, ParentProps, useContext } from "solid-js";
import CTX from "./ctx";
import { createFolder } from "@sv/apis/media";
import { VALIDATIONS } from "@sv/utils/validate";
import { COMMON_ERR_KEYS } from "@sv/utils/errors";

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

  const handleCreateFolderNameChange = (name: string) => {
    setCreateFolderName(name);
    isInvalidFolderName();
  };

  const isInvalidFolderName = (): boolean => {
    if (!createFolderName() || createFolderName().trim() === "") {
      setError("Folder name is required");
      return true;
    }
    if (createFolderName().length > VALIDATIONS.MAX_LENGTH) {
      setError("Folder name is too long");
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
      await createFolder(parentFolderId, createFolderName().trim());
      setCreateFolderName(""); // Reset the name after successful creation
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
      setIsCreateFolderModalOpen(false);
    }
  };

  const val = {
    isCreateFolderModalOpen,
    setIsCreateFolderModalOpen,
    createFolderName,
    isCreating,
    error,
    handleCreateFolderNameChange,
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
