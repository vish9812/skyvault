import { Accessor, createContext, Setter } from "solid-js";
import { FileInfo } from "@sv/apis/media/models";

export interface FileUploadState {
  file: File;
  progress: number;
  status: "pending" | "uploading" | "success" | "error";
  error?: string;
  result?: FileInfo;
}

interface CtxType {
  // Create folder related
  isCreateFolderModalOpen: Accessor<boolean>;
  setIsCreateFolderModalOpen: Setter<boolean>;
  createFolderName: Accessor<string>;
  isCreating: Accessor<boolean>;
  error: Accessor<string>;
  handleCreateFolderNameChange: (name: string) => void;
  handleCreateFolder: (parentFolderId?: string) => Promise<void>;

  // File upload related
  isFileUploadModalOpen: Accessor<boolean>;
  setIsFileUploadModalOpen: Setter<boolean>;
  fileUploads: Accessor<FileUploadState[]>;
  isUploading: Accessor<boolean>;
  uploadError: Accessor<string>;
  handleFileSelect: (files: FileList | null) => void;
  handleUploadFiles: (folderId?: string) => Promise<void>;
  clearUploads: () => void;
}

const CTX = createContext<CtxType>();

export default CTX;
