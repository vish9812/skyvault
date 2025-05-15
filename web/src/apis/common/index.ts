import { LOCAL_STORAGE_KEYS } from "@sv/utils/consts";

export const ROOT_URL = "http://localhost:8090/api/v1";
export const ROOT_URL_PUB = `${ROOT_URL}/pub`;

export function postJSONPub(url: string, data: any) {
  return fetch(`${ROOT_URL_PUB}/${url}`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(data),
  });
}

export function postJSON(url: string, data: any) {
  const token = localStorage.getItem(LOCAL_STORAGE_KEYS.TOKEN);
  if (!token) throw new Error("No token found");

  return fetch(`${ROOT_URL}/${url}`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify(data),
  });
}

export function getJSON(url: string) {
  const token = localStorage.getItem(LOCAL_STORAGE_KEYS.TOKEN);
  if (!token) throw new Error("No token found");

  return fetch(`${ROOT_URL}/${url}`, {
    headers: {
      Authorization: `Bearer ${token}`,
      Accept: "application/json",
    },
  });
}

export async function handleJSONResponse(res: Response) {
  if (!res.ok) {
    let msg = "";
    try {
      const data = await res.json();
      msg = data.code;
    } catch (e) {
      // If the response is not JSON, try to parse it as text
      if (e instanceof SyntaxError) {
        msg = await res.text();
      }
    }
    throw new Error(msg);
  }

  return res.json();
}
