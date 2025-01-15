import { clsx, type ClassValue } from "clsx";
import { twMerge } from "tailwind-merge";
import consts from "./consts";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

function convertFileToUrl(file: File) {
  return URL.createObjectURL(file);
}

function prettySize(sizeBytes: number) {
  if (sizeBytes < 1024) {
    return sizeBytes + " Bytes"; // Less than 1 KB, show in Bytes
  } else if (sizeBytes < 1024 * 1024) {
    const sizeInKB = sizeBytes / 1024;
    return sizeInKB.toFixed(1) + " KB"; // Less than 1 MB, show in KB
  } else if (sizeBytes < 1024 * 1024 * 1024) {
    const sizeInMB = sizeBytes / (1024 * 1024);
    return sizeInMB.toFixed(1) + " MB"; // Less than 1 GB, show in MB
  } else {
    const sizeInGB = sizeBytes / (1024 * 1024 * 1024);
    return sizeInGB.toFixed(1) + " GB"; // 1 GB or more, show in GB
  }
}

function formattedDateTime(isoString: string | null | undefined) {
  if (!isoString) return "-";

  const date = new Date(isoString);

  // Get hours and adjust for 12-hour format
  let hours = date.getHours();
  const minutes = date.getMinutes();
  const period = hours >= 12 ? "pm" : "am";

  // Convert hours to 12-hour format
  hours = hours % 12 || 12;

  // Format the time and date parts
  const time = `${hours}:${minutes.toString().padStart(2, "0")}${period}`;
  const day = date.getDate();
  const monthNames = [
    "Jan",
    "Feb",
    "Mar",
    "Apr",
    "May",
    "Jun",
    "Jul",
    "Aug",
    "Sep",
    "Oct",
    "Nov",
    "Dec",
  ];
  const month = monthNames[date.getMonth()];
  const year = date.getFullYear();
  const nowYear = new Date().getFullYear();
  const yearStr = year !== nowYear ? ` ${year}` : "";

  return `${time}, ${day} ${month}${yearStr}`;
}

// function prettyTimeAgo(date: Date) {
//   const userLocale = Intl.DateTimeFormat().resolvedOptions().locale;
//   const timeAgo = new TimeAgo(userLocale || "en-IN");
//   return timeAgo.format(date);
// }

function getFileIcon(extension: string | undefined, type: MediaType | string) {
  switch (extension) {
    // Text
    case "pdf":
      return "/assets/icons/file-pdf.svg";
    case "doc":
      return "/assets/icons/file-doc.svg";
    case "docx":
      return "/assets/icons/file-docx.svg";
    case "csv":
      return "/assets/icons/file-csv.svg";
    case "txt":
      return "/assets/icons/file-txt.svg";
    case "xls":
    case "xlsx":
      return "/assets/icons/file-document.svg";
    // Image
    case "svg":
      return "/assets/icons/file-image.svg";
    // Video
    case "mkv":
    case "mov":
    case "avi":
    case "wmv":
    case "mp4":
    case "flv":
    case "webm":
    case "m4v":
    case "3gp":
      return "/assets/icons/file-video.svg";
    // Audio
    case "mp3":
    case "mpeg":
    case "wav":
    case "aac":
    case "flac":
    case "ogg":
    case "wma":
    case "m4a":
    case "aiff":
    case "alac":
      return "/assets/icons/file-audio.svg";

    default:
      switch (type) {
        case consts.fileType.image:
          return "/assets/icons/file-image.svg";
        case consts.fileType.document:
          return "/assets/icons/file-document.svg";
        case consts.fileType.video:
          return "/assets/icons/file-video.svg";
        case consts.fileType.audio:
          return "/assets/icons/file-audio.svg";
        default:
          return "/assets/icons/file-other.svg";
      }
  }
}

function getFileTypeFromMimeType(mimeType: string) {
  const baseMime = mimeType.split("/")[0];
  switch (baseMime) {
    case "image":
      return consts.fileType.image;
    case "text":
      return consts.fileType.document;
    case "video":
      return consts.fileType.video;
    case "audio":
      return consts.fileType.audio;
    default:
      return "other";
  }
}

function getFileTypeAndExtension(fileName: string) {
  const extension = fileName.split(".").pop()?.toLowerCase();

  if (!extension) return { type: "other", extension: "" };

  const documentExtensions = [
    "pdf",
    "doc",
    "docx",
    "txt",
    "xls",
    "xlsx",
    "csv",
    "rtf",
    "ods",
    "ppt",
    "odp",
    "md",
    "html",
    "htm",
    "epub",
    "pages",
    "fig",
    "psd",
    "ai",
    "indd",
    "xd",
    "sketch",
    "afdesign",
    "afphoto",
    "afphoto",
  ];
  const imageExtensions = ["jpg", "jpeg", "png", "gif", "bmp", "svg", "webp"];
  const videoExtensions = ["mp4", "avi", "mov", "mkv", "webm"];
  const audioExtensions = ["mp3", "wav", "ogg", "flac"];

  if (documentExtensions.includes(extension))
    return { type: "document", extension };
  if (imageExtensions.includes(extension)) return { type: "image", extension };
  if (videoExtensions.includes(extension)) return { type: "video", extension };
  if (audioExtensions.includes(extension)) return { type: "audio", extension };

  return { type: "other", extension };
}

const utils = {
  getFileIcon,
  getFileTypeAndExtension,
  getFileTypeFromMimeType,
  convertFileToUrl,
  prettySize,
  formattedDateTime,
};

export default utils;
