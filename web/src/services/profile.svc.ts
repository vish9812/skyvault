import consts from "@/lib/consts";
import Profile from "@/profile/profile";

const profileURLPub = consts.configs.baseAPIPub + "/profile";

async function get(id: number): Promise<Profile> {
  const res = await fetch(profileURLPub + "/" + id, {
    method: "GET",
    headers: consts.headers.authJson(),
  });
  if (res.ok) {
    return res.json();
  } else {
    throw new Error(await res.text());
  }
}

const profileSvc = {
  get,
};

export default profileSvc;
