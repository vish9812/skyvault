import { AppCtxProvider } from "@sv/store/appCtxProvider";
import { ParentProps } from "solid-js";
import AppHeader from "./appHeader";
import AppNavigation from "./appNavigation";

function AppLayoutWithCtx(props: ParentProps) {
  return (
    <div class="bg-bg-subtle min-h-screen">
      <AppNavigation />
      <AppHeader />
      {/* Main Content */}
      <div class="md:ml-64 pl-2 pr-2 pt-16 pb-16 min-h-screen">
        <main>{props.children}</main>
      </div>
    </div>
  );
}

function AppLayout(props: ParentProps) {
  return (
    <AppCtxProvider>
      <AppLayoutWithCtx {...props} />
    </AppCtxProvider>
  );
}

export default AppLayout;
