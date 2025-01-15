import consts from "@/lib/consts";
import type { FileModel } from "@/lib/models";
import utils from "@/lib/utils";

const filesURLPvt = consts.configs.baseAPIPvt + "/media";

async function uploadFile(file: File, folderID: number | null): Promise<FileModel> {
  const formData = new FormData();
  formData.append("file", file);
  const queryParam = folderID ? `?folder-id=${folderID}` : "";

  const res = await fetch(filesURLPvt + queryParam, {
    method: "POST",
    headers: consts.headers.auth(),
    body: formData,
  });

  if (res.ok) {
    return res.json();
  } else {
    const errorText = await res.text();
    throw new Error(errorText);
  }
}

async function getFiles(folderID: number | null): Promise<FileModel[]> {
  const queryParam = folderID ? `?folder-id=${folderID}` : "";
  const res = await fetch(filesURLPvt + queryParam, {
    method: "GET",
    headers: consts.headers.auth(),
  });

  if (!res.ok) {
    const errorText = await res.text();
    throw new Error(errorText);
  }

  const files: FileModel[] = await res.json();
  for (const f of files) {
    const baseType = f.mimeType.split("/")[0];
    f.type = utils.getFileTypeFromMimeType(baseType);

    if (f.type === consts.fileType.image) {
      f.url = await getBlobDataURL(f.id);
    }
  }
  return files;
}

async function getBlob(fileID: number): Promise<Blob> {
  const res = await fetch(filesURLPvt + "/blob/" + fileID, {
    method: "GET",
    headers: consts.headers.auth(),
  });
  if (!res.ok) {
    const errorText = await res.text();
    throw new Error(errorText);
  }
  return res.blob();
}

async function getBlobDataURL(fileID: number): Promise<string> {
  const blob = await getBlob(fileID);
  return URL.createObjectURL(blob);
}

async function deleteFile(fileID: number): Promise<boolean> {
  const res = await fetch(filesURLPvt + "/" + fileID, {
    method: "DELETE",
    headers: consts.headers.auth(),
  });
  return res.ok;
}

const filesSvc = {
  uploadFile,
  getFiles,
  getBlob,
  getBlobDataURL,
  deleteFile,
};

export default filesSvc;
