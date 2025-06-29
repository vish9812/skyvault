export const COMMON_ERR_KEYS = {
  GENERIC: "COMMON_GENERIC_ERROR",
  INVALID: "COMMON_INVALID_VALUE",
  DUPLICATE: "COMMON_DUPLICATE_DATA",
  NO_DATA: "COMMON_NO_DATA",
} as const;

const errorMessages: Record<string, string> = {
  [COMMON_ERR_KEYS.GENERIC]: "Something went wrong. Please try again.",
  [COMMON_ERR_KEYS.INVALID]: "Please check your input and try again.",
  [COMMON_ERR_KEYS.DUPLICATE]: "This data already exists.",
  [COMMON_ERR_KEYS.NO_DATA]: "Data not found.",
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

export function defaultErrorMessage(code: string): string {
  return errorMessages[code] || errorMessages.COMMON_GENERIC_ERROR;
}

export function getAuthErrorMessage(code: string): string {
  return authErrorMessages[code] || defaultErrorMessage(code);
}

const fileUploadErrorMessages: Record<string, string> = {
  [COMMON_ERR_KEYS.DUPLICATE]: "This file already exists.",
  MEDIA_FILE_SIZE_LIMIT_EXCEEDED:
    "File size exceeds the maximum allowed size. Please try again with a smaller file.",
};

export function getFileUploadErrorMessage(code: string): string {
  return fileUploadErrorMessages[code] || defaultErrorMessage(code);
}
