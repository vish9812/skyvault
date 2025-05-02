import { createSignal, For, Show } from "solid-js";
import { z, ZodTypeAny } from "zod";
import { TextField } from "@kobalte/core/text-field";
import { Button } from "@kobalte/core/button";

interface AuthFormProps {
  schema: ZodTypeAny;
  onSubmit: (values: any) => void;
  submitLabel: string;
  error?: string;
  fields: ("fullName" | "email" | "password")[];
  showPasswordToggle?: boolean;
}

const fieldLabels: Record<string, string> = {
  fullName: "Full Name",
  email: "Email",
  password: "Password",
};

export default function AuthForm(props: AuthFormProps) {
  const [form, setForm] = createSignal<Record<string, string>>({});
  const [showPassword, setShowPassword] = createSignal(false);
  const [fieldErrors, setFieldErrors] = createSignal<Record<string, string>>(
    {}
  );

  const handleInput = (field: string, value: string) => {
    setForm((f) => ({ ...f, [field]: value }));
    setFieldErrors((e) => ({ ...e, [field]: "" }));
  };

  const handleSubmit = (e: Event) => {
    e.preventDefault();
    setFieldErrors({});
    const result = props.schema.safeParse(form());
    if (!result.success) {
      const errors: Record<string, string> = {};
      for (const err of result.error.errors) {
        if (err.path[0]) errors[err.path[0]] = err.message;
      }
      setFieldErrors(errors);
      return;
    }
    props.onSubmit(result.data);
  };

  return (
    <form
      class="w-full flex flex-col gap-4"
      onSubmit={handleSubmit}
      autocomplete="off"
    >
      <For each={props.fields}>
        {(field) => (
          <TextField
            name={field}
            value={form()[field] || ""}
            onChange={(value) => handleInput(field, value)}
            validationState={fieldErrors()[field] ? "invalid" : "valid"}
            required
          >
            <TextField.Label class="block text-sm font-medium text-gray-700 mb-1">
              {fieldLabels[field]}
            </TextField.Label>
            <Show
              when={field === "password" && props.showPasswordToggle}
              fallback={
                <TextField.Input
                  id={field}
                  type={field}
                  class="w-full px-3 py-2 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-primary"
                  autocomplete={
                    field === "password" ? "current-password" : "on"
                  }
                />
              }
            >
              <div class="relative">
                <TextField.Input
                  id={field}
                  type={showPassword() ? "text" : "password"}
                  class="w-full px-3 py-2 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-primary pr-10"
                  autocomplete="current-password"
                />
                <button
                  type="button"
                  class="absolute inset-y-0 right-2 flex items-center text-gray-400 hover:text-gray-600"
                  tabindex="-1"
                  onClick={() => setShowPassword((v) => !v)}
                  aria-label={
                    showPassword() ? "Hide password" : "Show password"
                  }
                >
                  {showPassword() ? (
                    <svg
                      xmlns="http://www.w3.org/2000/svg"
                      class="h-5 w-5"
                      fill="none"
                      viewBox="0 0 24 24"
                      stroke="currentColor"
                    >
                      <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        d="M13.875 18.825A10.05 10.05 0 0112 19c-5.523 0-10-4.477-10-10 0-1.657.403-3.22 1.125-4.575m1.875-2.25A9.956 9.956 0 0112 3c5.523 0 10 4.477 10 10 0 1.657-.403 3.22-1.125 4.575m-1.875 2.25A9.956 9.956 0 0112 21c-5.523 0-10-4.477-10-10 0-1.657.403-3.22 1.125-4.575"
                      />
                    </svg>
                  ) : (
                    <svg
                      xmlns="http://www.w3.org/2000/svg"
                      class="h-5 w-5"
                      fill="none"
                      viewBox="0 0 24 24"
                      stroke="currentColor"
                    >
                      <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        d="M15 12a3 3 0 11-6 0 3 3 0 016 0zm6 0c0 5.523-4.477 10-10 10S2 17.523 2 12 6.477 2 12 2s10 4.477 10 10z"
                      />
                    </svg>
                  )}
                </button>
              </div>
            </Show>
            <Show when={fieldErrors()[field]}>
              <TextField.ErrorMessage class="text-error text-xs mt-1">
                {fieldErrors()[field]}
              </TextField.ErrorMessage>
            </Show>
          </TextField>
        )}
      </For>
      <Show when={props.error}>
        <div class="text-error text-sm text-center mt-2">{props.error}</div>
      </Show>
      <Button
        type="submit"
        class="w-full bg-primary text-white font-semibold py-2 rounded hover:opacity-90 transition disabled:opacity-60"
        disabled={props.submitLabel.includes("...")}
      >
        {props.submitLabel}
      </Button>
    </form>
  );
}
