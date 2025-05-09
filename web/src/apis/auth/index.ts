import type {
  Profile,
  SignInReq,
  SignInRes,
  SignUpReq,
  SignUpRes,
} from "./models";
import { postJSONPub, handleJSONResponse } from "@sv/apis/common";

const BASE_URL = "auth";

export function getProfile(): Profile | null {
  const profile = localStorage.getItem("profile");
  if (!profile) return null;
  return JSON.parse(profile);
}

export function signOut() {
  localStorage.removeItem("token");
  localStorage.removeItem("profile");
}

async function handleAuthResponse(res: Response) {
  const data: SignInRes | SignUpRes = await handleJSONResponse(res);
  localStorage.setItem("token", data.token);
  localStorage.setItem("profile", JSON.stringify(data.profile));
}

export async function signIn({ email, password }: SignInReq): Promise<void> {
  const res = await postJSONPub(`${BASE_URL}/sign-in`, {
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
  const res = await postJSONPub(`${BASE_URL}/sign-up`, {
    fullName,
    email,
    password,
    provider: "email",
  });
  return handleAuthResponse(res);
}
