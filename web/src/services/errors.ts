export class ServerError extends Error {
  constructor(message: string) {
    super(message);
    this.name = "PublicError";
  }
}

// Common errors
const ErrCommonGeneric = "COMMON_GENERIC_ERROR";
const ErrCommonDuplicateData = "COMMON_DUPLICATE_DATA";
const ErrCommonNoData = "COMMON_NO_DATA";
const ErrCommonInvalidValue = "COMMON_INVALID_VALUE";
const ErrCommonNoAccess = "COMMON_NO_ACCESS";

// Media errors
const ErrMediaFileSizeLimitExceeded = "MEDIA_FILE_SIZE_LIMIT_EXCEEDED";

// Auth errors
const ErrAuthInvalidCredentials = "AUTH_INVALID_CREDENTIALS";
const ErrAuthInvalidToken = "AUTH_INVALID_TOKEN";
const ErrAuthTokenExpired = "AUTH_TOKEN_EXPIRED";
const ErrAuthInvalidProvider = "AUTH_INVALID_PROVIDER";

const errs = {
  ErrCommonGeneric,
  ErrCommonDuplicateData,
  ErrCommonNoData,
  ErrCommonInvalidValue,
  ErrCommonNoAccess,
  ErrMediaFileSizeLimitExceeded,
  ErrAuthInvalidCredentials,
  ErrAuthInvalidToken,
  ErrAuthTokenExpired,
  ErrAuthInvalidProvider,
};

export default errs;
