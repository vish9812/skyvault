import type {
  Profile,
  SignInReq,
  SignInRes,
  SignUpReq,
  SignUpRes,
} from "./models";
import { postJSONPub, handleJSONResponse } from "@sv/apis/common";
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

async function handleAuthResponse(res: Response) {
  const data: SignInRes | SignUpRes = await handleJSONResponse(res);
  localStorage.setItem(LOCAL_STORAGE_KEYS.TOKEN, data.token);
  localStorage.setItem(
    LOCAL_STORAGE_KEYS.PROFILE,
    JSON.stringify(data.profile)
  );
}

export async function signIn({ email, password }: SignInReq): Promise<void> {
  const res = await postJSONPub(`${urlAuth}/sign-in`, {
    provider: "email",
    providerUserId: email,
    password,
  });
  return handleAuthResponse(res);
}

export async function signUp({
  fullName,
  email,
  password,
}: SignUpReq): Promise<void> {
  const res = await postJSONPub(`${urlAuth}/sign-up`, {
    fullName,
    email,
    password,
    provider: "email",
  });
  return handleAuthResponse(res);
}
