export interface Profile {
  id: number;
  email: string;
  fullName: string;
  avatarBase64?: string;
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
