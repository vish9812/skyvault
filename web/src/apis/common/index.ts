import { LOCAL_STORAGE_KEYS } from "@sv/utils/consts";

export const ROOT_URL = `${import.meta.env.VITE_SERVER_URL}/api/v1`;
// public api root url
export const ROOT_URL_PUB = `${ROOT_URL}/pub`;

interface ErrRes {
  code: string;
}

export function postPub(url: string, data: any) {
  return fetch(`${ROOT_URL_PUB}/${url}`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(data),
  });
}

export function post(url: string, data: any) {
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

export function get(url: string) {
  const token = localStorage.getItem(LOCAL_STORAGE_KEYS.TOKEN);
  if (!token) throw new Error("No token found");

  return fetch(`${ROOT_URL}/${url}`, {
    headers: {
      Authorization: `Bearer ${token}`,
      Accept: "application/json",
    },
  });
}

export async function handleJSONResponse<T>(res: Response): Promise<T> {
  // Simulate a slow response
  return new Promise((resolve, reject) => {
    setTimeout(async () => {
      if (!res.ok) {
        let msg = "";
        try {
          const data = (await res.json()) as ErrRes;
          msg = data.code;
        } catch (e) {
          // If the response is not JSON, try to parse it as text
          if (e instanceof SyntaxError) {
            msg = await res.text();
          }
        }
        reject(new Error(msg));
      }

      resolve(res.json() as T);
    }, 700);
  });
}
