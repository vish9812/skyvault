import consts from "@/lib/consts";
import type { FileInfo } from "@/lib/models";
import utils from "@/lib/utils";
import { ServerError } from "./errors";

const mediaURLPvt = consts.configs.baseAPIPvt + "/media";

async function uploadFile(
  file: File,
  folderID: number | null
): Promise<FileInfo> {
  const formData = new FormData();
  formData.append("file", file);
  const queryParam = folderID ? `?folder-id=${folderID}` : "";

  const res = await fetch(mediaURLPvt + queryParam, {
    method: "POST",
    headers: consts.headers.auth(),
    body: formData,
  });

  const data = await res.json();

  if (res.ok) {
    return data;
  } else {
    throw new ServerError(data.code);
  }
}

async function getFilesInfo(folderID: number | null): Promise<FileInfo[]> {
  const queryParam = folderID ? `?folder-id=${folderID}` : "";
  const res = await fetch(mediaURLPvt + queryParam, {
    method: "GET",
    headers: consts.headers.auth(),
  });

  const data = await res.json();

  if (!res.ok) {
    throw new ServerError(data.code);
  }

  const infos: FileInfo[] = data.infos || [];
  for (const info of infos) {
    const baseType = info.mimeType.split("/")[0];
    info.type = utils.getFileTypeFromMimeType(baseType);

    if (info.type === consts.fileType.image) {
      const file = await getFile(info.id);
      info.url = utils.convertFileToUrl(file);
    }
  }
  return infos;
}

async function getFile(fileID: number): Promise<Blob> {
  const res = await fetch(mediaURLPvt + "/file/" + fileID, {
    method: "GET",
    headers: consts.headers.auth(),
  });
  if (!res.ok) {
    const err = await res.json();
    throw new ServerError(err.code);
  }
  return res.blob();
}

async function trashFile(fileID: number): Promise<void> {
  const res = await fetch(mediaURLPvt + "/" + fileID, {
    method: "DELETE",
    headers: consts.headers.auth(),
  });

  if (!res.ok) {
    const err = await res.json();
    throw new ServerError(err.code);
  }
}

const filesSvc = {
  uploadFile,
  getFilesInfo,
  getFile,
  trashFile,
};

export default filesSvc;
