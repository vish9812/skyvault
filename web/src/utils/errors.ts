const errorMessages: Record<string, string> = {
  COMMON_GENERIC_ERROR: "Something went wrong. Please try again.",
  COMMON_DUPLICATE_DATA: "This account already exists.",
  COMMON_NO_DATA: "No account found with these details.",
  COMMON_INVALID_VALUE: "Please check your input and try again.",
  AUTH_INVALID_CREDENTIALS: "Incorrect email or password.",
  AUTH_INVALID_TOKEN: "Your session is invalid. Please sign in again.",
  AUTH_TOKEN_EXPIRED: "Your session has expired. Please sign in again.",
  AUTH_WRONG_PROVIDER: "Please use the correct sign-in method.",
};

export function getErrorMessage(code?: string): string {
  if (!code) return errorMessages.COMMON_GENERIC_ERROR;
  return errorMessages[code] || errorMessages.COMMON_GENERIC_ERROR;
}
