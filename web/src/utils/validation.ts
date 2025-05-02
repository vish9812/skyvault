import { z } from "zod";

export const signInSchema = z.object({
  email: z
    .string()
    .email("Invalid email address")
    .max(255, "Email is too long"),
  password: z
    .string()
    .min(4, "Password must be at least 4 characters")
    .max(255, "Password is too long"),
});

export const signUpSchema = z.object({
  fullName: z
    .string()
    .min(1, "Full name is required")
    .max(255, "Full name is too long"),
  email: z
    .string()
    .email("Invalid email address")
    .max(255, "Email is too long"),
  password: z
    .string()
    .min(4, "Password must be at least 4 characters")
    .max(255, "Password is too long"),
});
