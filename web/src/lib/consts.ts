const configs = {
  baseAPIPvt: "/api/v1",
  baseAPIPub: "/api/v1/pub",
  maxFileSizeBytes: 5 * 1024 * 1024, // 5MB
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

export const navItems = [
  {
    name: "Home",
    icon: "/assets/icons/dashboard.svg",
    url: "/",
  },
  {
    name: "Documents",
    icon: "/assets/icons/documents.svg",
    url: "/documents",
  },
  {
    name: "Images",
    icon: "/assets/icons/images.svg",
    url: "/images",
  },
  {
    name: "Media",
    icon: "/assets/icons/video.svg",
    url: "/media",
  },
  {
    name: "Others",
    icon: "/assets/icons/others.svg",
    url: "/others",
  },
];

const consts = {
  configs,
  storageKeys,
  headers,
  pageRoutes,
  navItems,
};

export default consts;
