export interface FileInfo {
  id: string;
  ownerId: string;
  folderId?: string;
  name: string;
  size: number;
  extension?: string;
  mimeType: string;
  category: string;
  preview?: string;
  createdAt: string;
  updatedAt: string;
}

export interface FolderInfo {
  id: string;
  ownerId: string;
  parentFolderId?: string;
  name: string;
  createdAt: string;
  updatedAt: string;
}

export interface Page<T> {
  items: T[];
  prevCursor: string;
  nextCursor: string;
  hasMore: boolean;
}

export interface FolderContent {
  filePage: Page<FileInfo>;
  folderPage: Page<FolderInfo>;
}
