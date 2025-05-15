import { ParentComponent } from "solid-js";
import { useNavigate } from "@solidjs/router";
import { getProfile } from "@sv/apis/auth";
import { CLIENT_URLS } from "@sv/utils/consts";
import Header from "./header";
import Navigation from "./navigation";

const AppLayout: ParentComponent = (props) => {
  const navigate = useNavigate();
  const profile = getProfile();
  if (!profile) {
    navigate(CLIENT_URLS.SIGN_IN, { replace: true });
    return;
  }

  return (
    <div class="bg-bg-subtle min-h-screen">
      {/* Navigation */}
      <Navigation />

      <Header />

      {/* Main Content */}
      <div class="md:ml-64 pl-2 pr-2 pt-16 pb-16 min-h-screen">
        <main>{props.children}</main>
      </div>
    </div>
  );
};

export default AppLayout;
