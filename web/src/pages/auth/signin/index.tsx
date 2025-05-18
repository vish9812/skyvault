import useVM from "./useVM";
import { Link } from "@kobalte/core/link";
import { CLIENT_URLS } from "@sv/utils/consts";
import { TextField } from "@kobalte/core/text-field";
import { Show } from "solid-js";
import { Button } from "@kobalte/core/button";

function SignIn() {
  const vm = useVM();

  const hasFieldError = () => !!(vm.formErrors.email || vm.formErrors.password);
  const isDisabled = () => vm.loading() || hasFieldError();

  return (
    <>
      <h3>Sign In</h3>
      <form
        class="w-full flex flex-col gap-4"
        onSubmit={(e) => {
          e.preventDefault();
          vm.handleSubmit();
        }}
        autocomplete="off"
      >
        <TextField
          name="email"
          value={vm.formValues.email}
          onChange={(value) => vm.handleInput("email", value)}
          validationState={vm.formErrors.email ? "invalid" : "valid"}
        >
          <TextField.Label class="label">Email</TextField.Label>
          <TextField.Input
            id="email"
            type="email"
            autocomplete="email"
            class="input"
            classList={{
              "input-b-std": !vm.formErrors.email,
              "input-b-error": !!vm.formErrors.email,
            }}
          />
          <TextField.ErrorMessage class="input-t-error">
            {vm.formErrors.email}
          </TextField.ErrorMessage>
        </TextField>

        <TextField
          name="password"
          value={vm.formValues.password}
          onChange={(value) => vm.handleInput("password", value)}
          validationState={vm.formErrors.password ? "invalid" : "valid"}
        >
          <TextField.Label class="label">Password</TextField.Label>
          <div class="relative">
            <TextField.Input
              id="password"
              type={vm.showPassword() ? "text" : "password"}
              class="input"
              classList={{
                "input-b-std": !vm.formErrors.password,
                "input-b-error": !!vm.formErrors.password,
              }}
              autocomplete="current-password"
            />
            <Button
              class="absolute inset-y-0 right-2 flex-center text-neutral-lighter hover:text-neutral-light cursor-pointer"
              tabindex="-1"
              onClick={() => vm.setShowPassword((v) => !v)}
              aria-label={vm.showPassword() ? "Hide password" : "Show password"}
            >
              {vm.showPassword() ? (
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
            {vm.formErrors.password}
          </TextField.ErrorMessage>
        </TextField>

        <Show when={vm.apiError()}>
          <div class="text-error text-sm text-center mt-2">{vm.apiError()}</div>
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
          {vm.loading() ? "Signing In..." : "Sign In"}
        </Button>
      </form>

      <div class="mt-4">
        <span>Don't have an account? </span>
        <Link href={CLIENT_URLS.SIGN_UP} class="link font-semibold">
          Sign Up
        </Link>
      </div>
    </>
  );
}

export default SignIn;
