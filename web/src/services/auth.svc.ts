import consts from "@/lib/consts";
import Profile from "@/profile/profile";

const authURL = consts.configs.baseAPI + "/auth";

function isAuthenticated() {
  return !!localStorage.getItem(consts.storageKeys.auth_token);
}

interface SignUpReq {
  fullName: string;
  email: string;
  password: string;
}

interface SignInRes {
  profile: Profile;
  token: string;
}

async function handleAuthResponse(res: Response): Promise<Profile> {
  if (res.ok) {
    const data: SignInRes = await res.json();
    localStorage.setItem(consts.storageKeys.auth_token, data.token);
    localStorage.setItem(
      consts.storageKeys.profile,
      JSON.stringify(data.profile)
    );
    return data.profile;
  } else {
    const errorText = await res.text();
    throw new Error(errorText);
  }
}

async function signUp(req: SignUpReq): Promise<Profile> {
  const res = await fetch(authURL + "/sign-up", {
    method: "POST",
    headers: consts.headers.json,
    body: JSON.stringify(req),
  });
  return handleAuthResponse(res);
}

interface SignInReq {
  email: string;
  password: string;
}

async function signIn(req: SignInReq): Promise<Profile> {
  const res = await fetch(authURL + "/sign-in", {
    method: "POST",
    headers: consts.headers.json,
    body: JSON.stringify(req),
  });
  return handleAuthResponse(res);
}

function signOut() {
  localStorage.removeItem(consts.storageKeys.profile);
}

const authSvc = {
  isAuthenticated,
  signUp,
  signIn,
  signOut,
};

export default authSvc;
