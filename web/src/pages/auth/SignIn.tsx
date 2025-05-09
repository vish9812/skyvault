import { createSignal } from "solid-js";
import { useNavigate } from "@solidjs/router";
import { z } from "zod";
import AuthForm from "@sv/components/auth/AuthForm";
import { getErrorMessage } from "@sv/utils/errors";
import { signInSchema } from "@sv/pages/auth/validation";
import { signIn } from "../../apis/auth";
import type { SignInReq } from "../../apis/auth/models";

export default function SignIn() {
  const navigate = useNavigate();
  const [error, setError] = createSignal("");
  const [loading, setLoading] = createSignal(false);

  const handleSubmit = async (values: SignInReq) => {
    setLoading(true);
    setError("");
    try {
      const res = await signIn(values);
      localStorage.setItem("token", res.token);
      navigate("/");
    } catch (err: any) {
      setError(getErrorMessage(err?.code));
    } finally {
      setLoading(false);
    }
  };

  return (
    <>
      <AuthForm
        schema={signInSchema}
        onSubmit={handleSubmit}
        submitLabel={loading() ? "Signing In..." : "Sign In"}
        error={error()}
        fields={["email", "password"]}
        showPasswordToggle
      />
      <div class="mt-4 text-center">
        <span>Don't have an account? </span>
        <a href="/sign-up" class="text-primary font-semibold hover:underline">
          Sign Up
        </a>
      </div>
    </>
  );
}
