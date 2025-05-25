import type {
  Profile,
  SignInReq,
  SignInRes,
  SignUpReq,
  SignUpRes,
} from "./models";
import { postPub, handleJSONResponse } from "@sv/apis/common";
import { LOCAL_STORAGE_KEYS } from "@sv/utils/consts";

const urlAuth = "auth";

export function getProfile(): Profile | null {
  const profile = localStorage.getItem(LOCAL_STORAGE_KEYS.PROFILE);
  if (!profile) return null;
  return JSON.parse(profile);
}

export function signOut() {
  localStorage.removeItem(LOCAL_STORAGE_KEYS.TOKEN);
  localStorage.removeItem(LOCAL_STORAGE_KEYS.PROFILE);
}

async function handleAuthResponse<T extends SignInRes | SignUpRes>(
  res: Response
): Promise<void> {
  const data = await handleJSONResponse<T>(res);

  // TODO: Temporary: Implement profile preferences on server side
  data.profile.preferences = {
    contentView: "list",
  };

  localStorage.setItem(LOCAL_STORAGE_KEYS.TOKEN, data.token);
  localStorage.setItem(
    LOCAL_STORAGE_KEYS.PROFILE,
    JSON.stringify(data.profile)
  );
}

export async function signIn({ email, password }: SignInReq): Promise<void> {
  const res = await postPub(`${urlAuth}/sign-in`, {
    provider: "email",
    providerUserId: email,
    password,
  });
  return handleAuthResponse<SignInRes>(res);
}

export async function signUp({
  fullName,
  email,
  password,
}: SignUpReq): Promise<void> {
  const res = await postPub(`${urlAuth}/sign-up`, {
    fullName,
    email,
    password,
    provider: "email",
  });
  return handleAuthResponse<SignUpRes>(res);
}
