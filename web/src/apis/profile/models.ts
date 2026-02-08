export interface Profile {
  id: string;
  email: string;
  fullName: string;
  avatarBase64?: string;
  preferences: Preferences; //TODO: Make preferences json column in DB
}

export interface Preferences {
  contentView: "list" | "grid";
}

export interface StorageUsage {
  used: number;
  quota: number;
}
