import { Accessor, createContext, Setter } from "solid-js";

interface CtxType {
  isFolderModalOpen: Accessor<boolean>;
  setIsFolderModalOpen: Setter<boolean>;
  folderName: Accessor<string>;
  setFolderName: Setter<string>;
  isCreating: Accessor<boolean>;
  error: Accessor<string | null>;
  handleCreateFolder: (parentFolderId?: string) => Promise<void>;
}

const CTX = createContext<CtxType>();

export default CTX;
