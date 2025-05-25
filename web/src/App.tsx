import { Show, type Component } from "solid-js";
import AppRoutes from "./routes";
import { QueryClient, QueryClientProvider } from "@tanstack/solid-query";
import { SolidQueryDevtools } from "@tanstack/solid-query-devtools";
import { MetaProvider } from "@solidjs/meta";

const queryClient = new QueryClient();

const App: Component = () => {
  return (
    <MetaProvider>
      <QueryClientProvider client={queryClient}>
        <Show when={import.meta.env.DEV}>
          <SolidQueryDevtools />
        </Show>
        <AppRoutes />
      </QueryClientProvider>
    </MetaProvider>
  );
};

export default App;
