export interface FileInfo {
  id: number;
  ownerId: number;
  folderId?: number;
  name: string;
  size: number;
  extension?: string;
  mimeType: string;
  type: string;
  url: string;
  createdAt: string;
  updatedAt: string;
  trashedAt?: string;
}
