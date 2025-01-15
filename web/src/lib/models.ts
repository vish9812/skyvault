export interface FileModel {
  id: number;
  folderId?: number;
  ownerId: number;
  name: string;
  sizeBytes: number;
  mimeType: string;
  type: string;
  extension?: string;
  url: string;
  createdAt: string;
  updatedAt: string;
  trashedAt?: string;
}
