import { VALIDATIONS } from "@sv/utils/consts";
import { z } from "zod";

export const signInSchema = z.object({
  email: z
    .string()
    .email("Invalid email address")
    .max(VALIDATIONS.MAX_LENGTH, "Email is too long"),
  password: z
    .string()
    .min(4, "Password must be at least 4 characters")
    .max(VALIDATIONS.MAX_LENGTH, "Password is too long"),
});

export const signUpSchema = z.object({
  fullName: z
    .string()
    .min(1, "Full name is required")
    .max(VALIDATIONS.MAX_LENGTH, "Full name is too long"),
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
});
