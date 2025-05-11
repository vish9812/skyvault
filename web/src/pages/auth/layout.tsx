import { ParentProps } from "solid-js";
import { useNavigate } from "@solidjs/router";
import { getProfile } from "@sv/apis/auth";

function AuthLayout(props: ParentProps) {
  const navigate = useNavigate();
  const profile = getProfile();
  if (profile) {
    navigate("/", { replace: true });
  }
  return (
    <div class="min-h-screen flex-center gradient">
      <div class="w-full max-w-md bg-bg-subtle rounded-xl shadow-lg p-8 flex-center flex-col">
        <h1 class="gradient-text mb-2">SkyVault</h1>
        {props.children}
      </div>
    </div>
  );
}

export default AuthLayout;
