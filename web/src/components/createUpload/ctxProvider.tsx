import { createEffect, createSignal, ParentProps, useContext } from "solid-js";
import CTX, { FileUploadState } from "./ctx";
import { createFolder, uploadFiles } from "@sv/apis/media";
import { VALIDATIONS } from "@sv/utils/validate";
import { COMMON_ERR_KEYS } from "@sv/utils/errors";

export function CtxProvider(props: ParentProps) {
  // Create folder related states
  const [isCreateFolderModalOpen, setIsCreateFolderModalOpen] =
    createSignal(false);
  const [createFolderName, setCreateFolderName] = createSignal("");
  const [isCreating, setIsCreating] = createSignal(false);
  const [error, setError] = createSignal("");

  // File upload related states
  const [isFileUploadModalOpen, setIsFileUploadModalOpen] = createSignal(false);
  const [fileUploads, setFileUploads] = createSignal<FileUploadState[]>([]);
  const [isUploading, setIsUploading] = createSignal(false);
  const [uploadError, setUploadError] = createSignal("");
  const [currentFolderId, setCurrentFolderId] = createSignal<
    string | undefined
  >(undefined);

  // Reset the folder name when the modal is closed
  createEffect(() => {
    if (!isCreateFolderModalOpen()) {
      setCreateFolderName("");
      setError("");
    }
  });

  // Reset file uploads when modal is closed
  createEffect(() => {
    if (!isFileUploadModalOpen() && !isUploading()) {
      clearUploads();
    }
  });

  // Create folder functions
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

  // File upload functions
  const handleFileSelect = (files: FileList | null) => {
    if (!files || files.length === 0) return;

    setUploadError("");
    const newUploads: FileUploadState[] = Array.from(files).map((file) => ({
      file,
      progress: 0,
      status: "pending",
    }));

    setFileUploads(newUploads);
  };

  const handleUploadFiles = async (folderId?: string) => {
    if (fileUploads().length === 0 || isUploading()) return;

    setIsUploading(true);
    setUploadError("");
    const targetFolderId = folderId || currentFolderId();

    try {
      // Mark all files as uploading
      setFileUploads((prev) =>
        prev.map((item) => ({ ...item, status: "uploading" as const }))
      );

      // Convert to simple File array
      const files = fileUploads().map((upload) => upload.file);

      await uploadFiles(files, targetFolderId, (fileIndex, progress) => {
        setFileUploads((prev) => {
          const updated = [...prev];
          if (updated[fileIndex]) {
            updated[fileIndex] = {
              ...updated[fileIndex],
              progress,
            };
          }
          return updated;
        });
      });

      // Mark all as success
      setFileUploads((prev) =>
        prev.map((item) => ({ ...item, status: "success" as const }))
      );

      // Close the modal after a brief delay to show success state
      setTimeout(() => {
        setIsFileUploadModalOpen(false);
      }, 2000);
    } catch (err) {
      setUploadError(
        err instanceof Error ? err.message : "Failed to upload files"
      );

      // Mark remaining files as failed
      setFileUploads((prev) =>
        prev.map((item) =>
          item.status === "uploading"
            ? { ...item, status: "error" as const, error: "Upload failed" }
            : item
        )
      );
    } finally {
      setIsUploading(false);
    }
  };

  const clearUploads = () => {
    if (!isUploading()) {
      setFileUploads([]);
      setUploadError("");
    }
  };

  const setCurrentFolder = (folderId?: string) => {
    setCurrentFolderId(folderId);
  };

  const val = {
    // Create folder related values
    isCreateFolderModalOpen,
    setIsCreateFolderModalOpen,
    createFolderName,
    isCreating,
    error,
    handleCreateFolderNameChange,
    handleCreateFolder,

    // File upload related values
    isFileUploadModalOpen,
    setIsFileUploadModalOpen,
    fileUploads,
    isUploading,
    uploadError,
    handleFileSelect,
    handleUploadFiles,
    clearUploads,
    setCurrentFolder,
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
