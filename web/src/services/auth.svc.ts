import Config from "@/lib/config";
import Profile from "@/profile/profile";

const authURL = Config.API_URL + "/auth";

interface SignUpReq {
  fullName: string;
  email: string;
  password: string;
}

function signUp(req: SignUpReq): Promise<Profile> {
  return fetch(authURL + "/sign-up", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(req),
  }).then((response) => response.json());
}

interface SignInReq {
  email: string;
  password: string;
}
function signIn(req: SignInReq): Promise<Profile> {
  return fetch(authURL + "/sign-in", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(req),
  }).then((response) => response.json());
}

const authSvc = {
  signUp,
  signIn,
};

export default authSvc;
