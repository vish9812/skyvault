export enum CATEGORY {
  IMAGE = "image",
  VIDEO = "video",
  AUDIO = "audio",
  TEXT = "text",
  OTHER = "other",
}

export enum MIME_BASE {
  IMAGE = "image",
  VIDEO = "video",
  AUDIO = "audio",
  TEXT = "text",
}

/**
 * Convert MIME type to category
 * @param mimeType - Either full MIME type or base MIME type (e.g. "image/png" or "image")
 * @returns Category
 */
function mimeToCategory(mimeType: string): CATEGORY {
  const type = mimeType.split("/")[0];
  switch (type) {
    case MIME_BASE.IMAGE:
      return CATEGORY.IMAGE;
    case MIME_BASE.VIDEO:
      return CATEGORY.VIDEO;
    case MIME_BASE.AUDIO:
      return CATEGORY.AUDIO;
    case MIME_BASE.TEXT:
      return CATEGORY.TEXT;
    default:
      return CATEGORY.OTHER;
  }
}

function downloadBlob(blob: Blob, fileName: string) {
  const link = document.createElement("a");
  const url = window.URL.createObjectURL(blob);

  link.href = url;
  link.download = fileName;
  link.setAttribute("style", "display: none");

  document.body.appendChild(link);
  link.click();

  // Clean up
  document.body.removeChild(link);
  window.URL.revokeObjectURL(url);
}

const FileUtils = {
  mimeToCategory,
  downloadBlob,
} as const;

export default FileUtils;
