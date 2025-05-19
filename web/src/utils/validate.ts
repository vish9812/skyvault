import { z } from "zod";

export const VALIDATIONS = {
  MAX_LENGTH: 255,
  PASSWORD_MIN_LENGTH: 4,
} as const;

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

const Validators = {
  name,
  z: zSchemas,
} as const;

export default Validators;
