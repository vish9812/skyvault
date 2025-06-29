import { z } from "zod";

export const VALIDATIONS = {
  MAX_LENGTH: 255,
  PASSWORD_MIN_LENGTH: 4,
} as const;

/**
 * Validate any generic name in the app
 * @param name - Name to validate
 * @returns True if name is valid, false otherwise
 */
function name(name: string): boolean {
  return !!(
    name &&
    name.trim() !== "" &&
    name.length <= VALIDATIONS.MAX_LENGTH
  );
}

const zSchemas = {
  name: z
    .string()
    .min(1, "Name is required")
    .max(VALIDATIONS.MAX_LENGTH, "Name is too long"),
  email: z
    .string()
    .email("Invalid email address")
    .max(VALIDATIONS.MAX_LENGTH, "Email is too long"),
  password: z
    .string()
    .min(
      VALIDATIONS.PASSWORD_MIN_LENGTH,
      "Password must be at least 4 characters"
    )
    .max(VALIDATIONS.MAX_LENGTH, "Password is too long"),
} as const;

const Validate = {
  name,
  z: zSchemas,
} as const;

export default Validate;
