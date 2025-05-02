import { createSignal } from "solid-js";
import { useNavigate } from "@solidjs/router";
import { z } from "zod";
import { Link } from "@kobalte/core/link";
import AuthForm from "@sv/components/auth/AuthForm";
import { getErrorMessage } from "@sv/utils/errors";
import { signUpSchema } from "@sv/utils/validation";
import { signUp } from "@sv/utils/api";

export default function SignUp() {
  console.log("SignUp");
  const navigate = useNavigate();
  const [error, setError] = createSignal("");
  const [loading, setLoading] = createSignal(false);

  const handleSubmit = async (values: z.infer<typeof signUpSchema>) => {
    setLoading(true);
    setError("");
    try {
      const res = await signUp(values);
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
