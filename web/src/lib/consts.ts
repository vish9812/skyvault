const configs = {
  baseAPI: "/api/v1",
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

const consts = {
  configs,
  storageKeys,
  headers,
  pageRoutes,
};

export default consts;
