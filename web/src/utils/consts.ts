export const ROOT_FOLDER_ID = "-1";
export const ROOT_FOLDER_NAME = "Root";

export const CLIENT_URLS = {
  SIGN_UP: "/auth/sign-up",
  SIGN_IN: "/auth/sign-in",
  DRIVE: "/drive",
  SHARED: "/shared",
  FAVORITES: "/favorites",
  TRASH: "/trash",
} as const;

export enum FILE_CATEGORIES {
  IMAGES = "images",
  VIDEOS = "videos",
  AUDIOS = "audios",
  DOCUMENTS = "documents",
  OTHERS = "others",
}

export const FOLDER_CONTENT_TYPES = {
  FOLDER: "folder",
  FILE: "file",
} as const;

export const CONTENT_VIEWS = {
  LIST: "list",
  GRID: "grid",
} as const;

export const LOCAL_STORAGE_KEYS = {
  CONTENT_VIEW: "content-view",
  PROFILE: "profile",
  TOKEN: "token",
} as const;
