import { Accessor, createContext, Setter } from "solid-js";

interface CtxType {
  isCreateFolderModalOpen: Accessor<boolean>;
  setIsCreateFolderModalOpen: Setter<boolean>;
  createFolderName: Accessor<string>;
  isCreating: Accessor<boolean>;
  error: Accessor<string>;
  handleCreateFolderNameChange: (name: string) => void;
  handleCreateFolder: (parentFolderId?: string) => Promise<void>;
}

const CTX = createContext<CtxType>();

export default CTX;
