import { useNavigate } from "@solidjs/router";
import { createSignal } from "solid-js";
import { getAuthErrorMessage } from "@sv/utils/errors";
import { signIn } from "@sv/apis/auth";
import { CLIENT_URLS, VALIDATIONS } from "@sv/utils/consts";
import { createStore, unwrap } from "solid-js/store";
import { z } from "zod";

export const schema = z.object({
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

type Fields = {
  email: string;
  password: string;
};

const emptyFields: Fields = {
  email: "",
  password: "",
} as const;

function useVM() {
  const [formValues, setFormValues] = createStore({ ...emptyFields });
  const [formErrors, setFormErrors] = createStore({ ...emptyFields });
  const [apiError, setApiError] = createSignal("");
  const [loading, setLoading] = createSignal(false);
  const [showPassword, setShowPassword] = createSignal(false);

  const navigate = useNavigate();

  const handleInput = (field: keyof Fields, value: string) => {
    setFormValues(field, value);
    const fieldSchema = schema.shape[field];
    const result = fieldSchema.safeParse(value);
    if (result.success) {
      setFormErrors(field, "");
    } else {
      setFormErrors(field, result.error.issues[0].message);
    }
  };

  const handleSubmit = async () => {
    setApiError("");
    setFormErrors(emptyFields);
    // Validate form
    const result = schema.safeParse(unwrap(formValues));
    if (!result.success) {
      for (const err of result.error.errors) {
        if (err.path[0]) {
          setFormErrors(err.path[0] as keyof typeof formErrors, err.message);
        }
      }
      return;
    }
    // Submit form
    setLoading(true);
    try {
      await signIn(unwrap(formValues));
      navigate(CLIENT_URLS.HOME);
    } catch (err: any) {
      setApiError(getAuthErrorMessage(err?.message));
    } finally {
      setLoading(false);
    }
  };

  return {
    formValues,
    formErrors,
    loading,
    showPassword,
    apiError,
    handleSubmit,
    handleInput,
    setShowPassword,
  };
}

export default useVM;
