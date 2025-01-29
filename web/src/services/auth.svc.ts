import consts from "@/lib/consts";
import Profile from "@/profile/profile";

const authURLPub = consts.configs.baseAPIPub + "/auth";

function profile(): Profile | null {
  const profileStr = localStorage.getItem(consts.storageKeys.profile);
  if (!profileStr) {
    return null;
  }

  return JSON.parse(profileStr);
}

interface SignUpReq {
  fullName: string;
  email: string;
  password: string;
  provider: string;
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
  const res = await fetch(authURLPub + "/sign-up", {
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
  const res = await fetch(authURLPub + "/sign-in", {
    method: "POST",
    headers: consts.headers.json,
    body: JSON.stringify(req),
  });
  return handleAuthResponse(res);
}

function signOut() {
  localStorage.removeItem(consts.storageKeys.profile);
  localStorage.removeItem(consts.storageKeys.auth_token);
  // Redirect to sign-in page
  window.location.href = consts.pageRoutes.signIn;
}

const authSvc = {
  profile,
  signUp,
  signIn,
  signOut,
};

export default authSvc;
