import type { CONTENT_VIEWS } from "@sv/utils/consts";

export interface Profile {
  id: string;
  email: string;
  fullName: string;
  avatarBase64?: string;
  preferences: Preferences; //TODO: Make preferences json column in DB
}

export type ContentView = (typeof CONTENT_VIEWS)[keyof typeof CONTENT_VIEWS];

export interface Preferences {
  contentView: ContentView;
}

export interface StorageUsage {
  used: number;
  quota: number;
}
