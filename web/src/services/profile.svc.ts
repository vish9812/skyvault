import consts from "@/lib/consts";
import Profile from "@/profile/profile";
import { ServerError } from "./errors";

const profileURLPub = consts.configs.baseAPIPvt + "/profile";

async function get(id: number): Promise<Profile> {
  const res = await fetch(profileURLPub + "/" + id, {
    method: "GET",
    headers: consts.headers.authJson(),
  });

  const data = await res.json();

  if (res.ok) {
    return data;
  } else {
    throw new ServerError(data.code);
  }
}

const profileSvc = {
  get,
};

export default profileSvc;
