import { CATEGORY } from "@sv/utils/fileUtils";

export interface UploadConfig {
  maxDirectUploadSizeMB: number;
  maxChunkSizeMB: number;
}

export interface UploadFileInfo {
  id: string;
  file: File;
  progress: number;
  status: "pending" | "uploading" | "success" | "error";
  error?: string;
}

export interface UploadFileResult {
  clientId: string;
  file: Promise<FileInfo>;
}

export interface FileInfo {
  id: string;
  ownerId: string;
  folderId?: string;
  name: string;
  size: number;
  extension?: string;
  mimeType: string;
  category: CATEGORY;
  previewBase64?: string;
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
