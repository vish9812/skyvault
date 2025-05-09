import useViewModel from "./useViewModel";
import { Link } from "@kobalte/core/link";
import { CLIENT_URLS } from "@sv/utils/consts";
import { TextField } from "@kobalte/core/text-field";
import { Show } from "solid-js";
import { Button } from "@kobalte/core/button";

function SignUp() {
  const {
    apiError,
    formErrors,
    loading,
    formValues,
    handleInput,
    handleSubmit,
    showPassword,
    setShowPassword,
  } = useViewModel();

  const hasFieldError = () =>
    !!(formErrors.fullName || formErrors.email || formErrors.password);

  const isDisabled = () => loading() || hasFieldError();

  return (
    <>
      <h3>Sign Up</h3>
      <form
        class="w-full flex flex-col gap-4"
        onSubmit={(e) => {
          e.preventDefault();
          handleSubmit();
        }}
        autocomplete="off"
      >
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
            class="input"
            classList={{
              "input-b-std": !formErrors.fullName,
              "input-b-error": !!formErrors.fullName,
            }}
          />
          <TextField.ErrorMessage class="input-t-error">
            {formErrors.fullName}
          </TextField.ErrorMessage>
        </TextField>

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
            class="input"
            classList={{
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
              class="input"
              classList={{
                "input-b-std": !formErrors.password,
                "input-b-error": !!formErrors.password,
              }}
              autocomplete="new-password"
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
          class="btn"
          classList={{
            "btn-disabled": isDisabled(),
            "btn-gradient": !isDisabled(),
          }}
          disabled={isDisabled()}
        >
          {loading() ? "Signing Up..." : "Sign Up"}
        </Button>
      </form>

      <div class="mt-4">
        <span>Already have an account? </span>
        <Link href={CLIENT_URLS.SIGN_IN} class="link font-semibold">
          Sign In
        </Link>
      </div>
    </>
  );
}

export default SignUp;
