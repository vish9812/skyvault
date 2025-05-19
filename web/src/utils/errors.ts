export const COMMON_ERR_KEYS = {
  GENERIC: "COMMON_GENERIC_ERROR",
  INVALID: "COMMON_INVALID_VALUE",
  DUPLICATE: "COMMON_DUPLICATE_DATA",
  NO_DATA: "COMMON_NO_DATA",
} as const;

const errorMessages: Record<string, string> = {
  [COMMON_ERR_KEYS.GENERIC]: "Something went wrong. Please try again.",
  [COMMON_ERR_KEYS.INVALID]: "Please check your input and try again.",
};

const authErrorMessages: Record<string, string> = {
  [COMMON_ERR_KEYS.DUPLICATE]: "This account already exists. Please sign in.",
  [COMMON_ERR_KEYS.NO_DATA]:
    "No account found with these details. Please sign up.",
  AUTH_INVALID_CREDENTIALS: "Incorrect email or password.",
  AUTH_INVALID_TOKEN: "Your session is invalid. Please sign in again.",
  AUTH_TOKEN_EXPIRED: "Your session has expired. Please sign in again.",
  AUTH_WRONG_PROVIDER: "Please use the correct sign-in method.",
};

export function getAuthErrorMessage(code?: string): string {
  if (!code) return errorMessages.COMMON_GENERIC_ERROR;

  return (
    authErrorMessages[code] ||
    errorMessages[code] ||
    errorMessages.COMMON_GENERIC_ERROR
  );
}
