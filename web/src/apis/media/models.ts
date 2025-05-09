export interface FileInfo {
  id: number;
  ownerId: number;
  folderId?: number;
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
  id: number;
  ownerId: number;
  parentFolderId?: number;
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
