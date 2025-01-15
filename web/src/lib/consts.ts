const configs = {
  baseAPIPvt: "/api/v1",
  baseAPIPub: "/api/v1/pub",
  maxFileSizeBytes: 100 * 1024 * 1024, // 100MB
} as const;

const fileType = {
  document: "document",
  image: "image",
  audio: "audio",
  video: "video",
  others: "others",
} as const;

const storageKeys = {
  auth_token: "auth_token",
  profile: "profile",
} as const;

const headers = {
  json: { "Content-Type": "application/json" } as const,
  auth: () => {
    const token = localStorage.getItem(storageKeys.auth_token);
    return { Authorization: `Bearer ${token}` };
  },
  authJson: () => ({ ...headers.json, ...headers.auth() }),
} as const;

const pageRoutes = {
  home: "/",
  signIn: "/sign-in",
  signUp: "/sign-up",
} as const;

const navItems = [
  {
    name: "Home",
    icon: "/assets/icons/dashboard.svg",
    url: "/",
  },
  {
    name: "Documents",
    icon: "/assets/icons/documents.svg",
    url: "/media/documents",
  },
  {
    name: "Images",
    icon: "/assets/icons/images.svg",
    url: "/media/images",
  },
  {
    name: "Media",
    icon: "/assets/icons/video.svg",
    url: "/media/videos",
  },
  {
    name: "Others",
    icon: "/assets/icons/others.svg",
    url: "/media/others",
  },
];

const actionsDropdownItems = [
  {
    label: "Rename",
    icon: "/assets/icons/edit.svg",
    value: "rename",
  },
  {
    label: "Details",
    icon: "/assets/icons/info.svg",
    value: "details",
  },
  {
    label: "Share",
    icon: "/assets/icons/share.svg",
    value: "share",
  },
  {
    label: "Download",
    icon: "/assets/icons/download.svg",
    value: "download",
  },
  {
    label: "Delete",
    icon: "/assets/icons/delete.svg",
    value: "delete",
  },
];

const consts = {
  configs,
  storageKeys,
  headers,
  pageRoutes,
  navItems,
  actionsDropdownItems,
  fileType,
};

export default consts;
