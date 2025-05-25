import { FILE_CATEGORIES } from "@sv/utils/consts";

export interface FileInfo {
  id: string;
  ownerId: string;
  folderId?: string;
  name: string;
  size: number;
  extension?: string;
  mimeType: string;
  category: FILE_CATEGORIES;
  preview?: string;
  createdAt: string;
  updatedAt: string;
}

export interface FolderInfo {
  id: string;
  ownerId: string;
  name: string;
  parentFolderId?: string;
  createdAt: string;
  updatedAt: string;
  ancestors: BaseInfo[];
}

export interface Page<T> {
  items: T[];
  prevCursor: string;
  nextCursor: string;
  hasMore: boolean;
}

export interface BaseInfo {
  id: string;
  name: string;
}

export interface FolderContent {
  filePage: Page<FileInfo>;
  folderPage: Page<FolderInfo>;
}
