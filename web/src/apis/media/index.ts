import type { FolderContent } from "./models";

export async function fetchRootContent(): Promise<FolderContent> {
  const token = localStorage.getItem("token");
  const res = await fetch(
    "http://localhost:8090/api/v1/media/folders/0/content",
    {
      headers: { Authorization: `Bearer ${token}` },
    }
  );
  if (!res.ok) throw new Error("Failed to fetch folder content");
  return res.json();
}
