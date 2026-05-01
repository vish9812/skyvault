import { useNavigate } from "@solidjs/router";
import { isLoggedIn } from "@sv/apis/auth";
import { CLIENT_URLS } from "@sv/utils/consts";
import type { ParentProps } from "solid-js";

function AuthLayout(props: ParentProps) {
  const navigate = useNavigate();
  if (isLoggedIn()) {
    navigate(CLIENT_URLS.DRIVE, { replace: true });
    return;
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
