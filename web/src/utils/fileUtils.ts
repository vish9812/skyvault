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

/**
 * Create image preview from file
 * @param file - File to create preview from
 * @returns Promise that resolves to the preview image. Rejects if file is not an image.
 */
function createImagePreview(file: File): Promise<string> {
  return new Promise((resolve, reject) => {
    if (!file.type.startsWith(MIME_BASE.IMAGE + "/")) {
      reject(new Error("Not an image file"));
      return;
    }

    const reader = new FileReader();
    reader.onload = (e) => {
      resolve(e.target?.result as string);
    };
    reader.onerror = reject;
    reader.readAsDataURL(file);
  });
}

const FileUtils = {
  mimeToCategory,
  createImagePreview,
} as const;

export default FileUtils;
