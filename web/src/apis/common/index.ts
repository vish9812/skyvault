import { LOCAL_STORAGE_KEYS } from "@sv/utils/consts";

// In production (same origin), use relative URL. In dev, use full URL to backend.
const SERVER_URL = import.meta.env.VITE_SERVER_URL || "";
export const ROOT_URL = `${SERVER_URL}/api/v1`;
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

export function postFormData(
  url: string,
  formData: FormData,
  onProgress?: (progress: number) => void
) {
  const token = localStorage.getItem(LOCAL_STORAGE_KEYS.TOKEN);
  if (!token) throw new Error("No token found");

  return new Promise<Response>((resolve, reject) => {
    const xhr = new XMLHttpRequest();
    xhr.open("POST", `${ROOT_URL}/${url}`);
    xhr.setRequestHeader("Authorization", `Bearer ${token}`);

    xhr.upload.onprogress = (event) => {
      if (event.lengthComputable && onProgress) {
        const progress = Math.round((event.loaded / event.total) * 100);
        onProgress(progress);
      }
    };

    xhr.onload = () => {
      const response = new Response(xhr.response, {
        status: xhr.status,
        statusText: xhr.statusText,
        headers: new Headers({
          "Content-Type": "application/json",
        }),
      });
      resolve(response);
    };

    xhr.onerror = () => {
      const msg = xhr.statusText || "SkyVault: Network Error";
      console.log(msg);
      reject(new Error(msg));
    };
    xhr.ontimeout = () => {
      const msg = xhr.statusText || "SkyVault: Request Timeout";
      console.log(msg);
      reject(new Error(msg));
    };

    xhr.send(formData);
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
    throw new Error(msg);
  }

  return res.json() as T;
}

export async function handleBlobResponse(res: Response): Promise<Blob> {
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
    throw new Error(msg);
  }
  return res.blob();
}
