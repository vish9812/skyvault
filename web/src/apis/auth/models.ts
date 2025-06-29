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

export interface SignInReq {
  email: string;
  password: string;
}

export interface SignInRes {
  token: string;
  profile: Profile;
}

export interface SignUpReq {
  fullName: string;
  email: string;
  password: string;
}

export interface SignUpRes {
  token: string;
  profile: Profile;
}
