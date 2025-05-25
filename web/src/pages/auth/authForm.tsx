import { Button } from "@kobalte/core/button";
import { TextField } from "@kobalte/core/text-field";
import { A, useNavigate } from "@solidjs/router";
import { signIn, signUp } from "@sv/apis/auth";
import { CLIENT_URLS } from "@sv/utils/consts";
import { getAuthErrorMessage } from "@sv/utils/errors";
import Validators from "@sv/utils/validate";
import { createSignal, Show } from "solid-js";
import { createStore, unwrap } from "solid-js/store";
import { z } from "zod";

const schema = z.object({
  fullName: Validators.z.name,
  email: Validators.z.email,
  password: Validators.z.password,
});

type FieldsType = z.infer<typeof schema>;

const emptyFields = (): FieldsType => ({
  fullName: "",
  email: "",
  password: "",
});

interface AuthFormProps {
  isSignUp: boolean;
}

function AuthForm(props: AuthFormProps) {
  const navigate = useNavigate();

  const [formValues, setFormValues] = createStore<FieldsType>(emptyFields());
  const [formErrors, setFormErrors] = createStore<FieldsType>(emptyFields());
  const [apiError, setApiError] = createSignal("");
  const [loading, setLoading] = createSignal(false);
  const [showPassword, setShowPassword] = createSignal(false);

  const isDisabled = () =>
    !!(
      loading() ||
      formErrors.email ||
      formErrors.password ||
      (props.isSignUp && formErrors.fullName)
    );

  const handleInput = (field: keyof FieldsType, value: string) => {
    setFormValues(field, value);

    const fieldSchema = schema.shape[field];
    const result = fieldSchema.safeParse(value);
    setFormErrors(field, result.success ? "" : result.error.issues[0].message);
  };

  const handleSubmit = async () => {
    console.log("handleSubmit");
    setLoading(true);
    setApiError("");
    setFormErrors(emptyFields());

    // Validate form
    const formSchema = props.isSignUp
      ? schema
      : schema.omit({ fullName: true });

    const result = formSchema.safeParse(unwrap(formValues));
    console.log("result", result);
    if (!result.success) {
      for (const err of result.error.errors) {
        if (err.path[0]) {
          setFormErrors(err.path[0] as keyof typeof formErrors, err.message);
        }
      }
      setLoading(false);
      return;
    }

    // Submit form
    try {
      console.log("formValues", formValues);
      const reqData = unwrap(formValues);
      console.log("reqData", reqData);
      await (props.isSignUp ? signUp(reqData) : signIn(reqData));
      navigate(CLIENT_URLS.DRIVE, { replace: true });
    } catch (err) {
      setApiError(getAuthErrorMessage((err as Error).message));
    } finally {
      setLoading(false);
    }
  };

  return (
    <>
      <h3>{props.isSignUp ? "Sign Up" : "Sign In"}</h3>
      <form
        class="w-full flex flex-col gap-4"
        onSubmit={(e) => {
          e.preventDefault();
          handleSubmit();
        }}
        autocomplete="on"
      >
        <Show when={props.isSignUp}>
          <TextField
            name="fullName"
            value={formValues.fullName}
            onChange={(value) => handleInput("fullName", value)}
            validationState={formErrors.fullName ? "invalid" : "valid"}
          >
            <TextField.Label class="label">Full Name</TextField.Label>
            <TextField.Input
              id="fullName"
              type="text"
              autocomplete="name"
              classList={{
                input: true,
                "input-b-std": !formErrors.fullName,
                "input-b-error": !!formErrors.fullName,
              }}
            />
            <TextField.ErrorMessage class="input-t-error">
              {formErrors.fullName}
            </TextField.ErrorMessage>
          </TextField>
        </Show>

        <TextField
          name="email"
          value={formValues.email}
          onChange={(value) => handleInput("email", value)}
          validationState={formErrors.email ? "invalid" : "valid"}
        >
          <TextField.Label class="label">Email</TextField.Label>
          <TextField.Input
            id="email"
            type="email"
            autocomplete="email"
            classList={{
              input: true,
              "input-b-std": !formErrors.email,
              "input-b-error": !!formErrors.email,
            }}
          />
          <TextField.ErrorMessage class="input-t-error">
            {formErrors.email}
          </TextField.ErrorMessage>
        </TextField>

        <TextField
          name="password"
          value={formValues.password}
          onChange={(value) => handleInput("password", value)}
          validationState={formErrors.password ? "invalid" : "valid"}
        >
          <TextField.Label class="label">Password</TextField.Label>
          <div class="relative">
            <TextField.Input
              id="password"
              type={showPassword() ? "text" : "password"}
              classList={{
                input: true,
                "input-b-std": !formErrors.password,
                "input-b-error": !!formErrors.password,
              }}
              autocomplete={
                props.isSignUp ? "new-password" : "current-password"
              }
            />
            <Button
              class="absolute inset-y-0 right-2 flex-center text-neutral-lighter hover:text-neutral-light cursor-pointer"
              tabindex="-1"
              onClick={() => setShowPassword((v) => !v)}
              aria-label={showPassword() ? "Hide password" : "Show password"}
            >
              {showPassword() ? (
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke-width="1.5"
                  stroke="currentColor"
                  class="size-6"
                >
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    d="M2.036 12.322a1.012 1.012 0 0 1 0-.639C3.423 7.51 7.36 4.5 12 4.5c4.638 0 8.573 3.007 9.963 7.178.07.207.07.431 0 .639C20.577 16.49 16.64 19.5 12 19.5c-4.638 0-8.573-3.007-9.963-7.178Z"
                  />
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    d="M15 12a3 3 0 1 1-6 0 3 3 0 0 1 6 0Z"
                  />
                </svg>
              ) : (
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke-width="1.5"
                  stroke="currentColor"
                  class="size-6"
                >
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    d="M3.98 8.223A10.477 10.477 0 0 0 1.934 12C3.226 16.338 7.244 19.5 12 19.5c.993 0 1.953-.138 2.863-.395M6.228 6.228A10.451 10.451 0 0 1 12 4.5c4.756 0 8.773 3.162 10.065 7.498a10.522 10.522 0 0 1-4.293 5.774M6.228 6.228 3 3m3.228 3.228 3.65 3.65m7.894 7.894L21 21m-3.228-3.228-3.65-3.65m0 0a3 3 0 1 0-4.243-4.243m4.242 4.242L9.88 9.88"
                  />
                </svg>
              )}
            </Button>
          </div>
          <TextField.ErrorMessage class="input-t-error">
            {formErrors.password}
          </TextField.ErrorMessage>
        </TextField>

        <Show when={apiError()}>
          <div class="text-error text-sm text-center mt-2">{apiError()}</div>
        </Show>

        <Button
          type="submit"
          classList={{
            btn: true,
            "btn-disabled": isDisabled(),
            "btn-gradient": !isDisabled(),
          }}
          disabled={isDisabled()}
        >
          {loading()
            ? props.isSignUp
              ? "Signing Up..."
              : "Signing In..."
            : props.isSignUp
            ? "Sign Up"
            : "Sign In"}
        </Button>
      </form>

      <div class="mt-4">
        <span>
          {props.isSignUp
            ? "Already have an account? "
            : "Don't have an account? "}
        </span>
        <A
          href={props.isSignUp ? CLIENT_URLS.SIGN_IN : CLIENT_URLS.SIGN_UP}
          class="link font-semibold"
        >
          {props.isSignUp ? "Sign In" : "Sign Up"}
        </A>
      </div>
    </>
  );
}

export default AuthForm;
