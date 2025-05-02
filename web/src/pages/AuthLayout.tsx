import type { ParentComponent } from "solid-js";
import { onMount } from "solid-js";
import { useNavigate } from "@solidjs/router";

interface AuthLayoutProps {
  title: string;
}

const AuthLayout: ParentComponent<AuthLayoutProps> = (props) => {
  const navigate = useNavigate();
  onMount(() => {
    const token = localStorage.getItem("token");
    if (token) {
      navigate("/", { replace: true });
    }
  });
  return (
    <div class="min-h-screen flex items-center justify-center bg-primary">
      <div class="w-full max-w-md bg-white/80 rounded-xl shadow-lg p-8 flex flex-col items-center">
        <h1 class="text-3xl font-extrabold bg-gradient-to-r from-indigo-500 to-sky-400 bg-clip-text text-transparent mb-2 select-none">
          SkyVault
        </h1>
        <h2 class="text-xl font-semibold text-gray-800 mb-6 text-center">
          {props.title}
        </h2>
        {props.children}
      </div>
    </div>
  );
};

export default AuthLayout;
