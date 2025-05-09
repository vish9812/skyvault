import { createSignal } from "solid-js";
import { useNavigate } from "@solidjs/router";
import { Link } from "@kobalte/core/link";
import AuthForm from "@sv/components/auth/AuthForm";
import { getErrorMessage } from "@sv/utils/errors";
import { signUpSchema } from "@sv/pages/auth/validation";
import { signUp } from "@sv/apis/auth";
import type { SignUpReq } from "@sv/apis/auth/models";

export default function SignUp() {
  const navigate = useNavigate();
  const [error, setError] = createSignal("");
  const [loading, setLoading] = createSignal(false);

  const handleSubmit = async (values: SignUpReq) => {
    setLoading(true);
    setError("");
    try {
      await signUp(values);
      navigate("/");
    } catch (err: any) {
      setError(getErrorMessage(err?.message));
    } finally {
      setLoading(false);
    }
  };

  return (
    <>
      <AuthForm
        schema={signUpSchema}
        onSubmit={handleSubmit}
        submitLabel={loading() ? "Signing Up..." : "Sign Up"}
        error={error()}
        fields={["fullName", "email", "password"]}
        showPasswordToggle
      />
      <div class="mt-4 text-center">
        <span>Already have an account? </span>
        <Link
          href="/sign-in"
          class="text-primary font-semibold hover:underline"
        >
          Sign In
        </Link>
      </div>
    </>
  );
}
