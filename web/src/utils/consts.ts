export const VALIDATIONS = {
  MAX_LENGTH: 255,
  PASSWORD_MIN_LENGTH: 4,
} as const;

export const CLIENT_URLS = {
  SIGN_UP: "/auth/sign-up",
  SIGN_IN: "/auth/sign-in",
  HOME: "/",
} as const;

export const FILE_CATEGORIES = {
  IMAGES: "images",
  VIDEOS: "videos",
  AUDIOS: "audios",
  DOCUMENTS: "documents",
  OTHERS: "others",
} as const;

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
