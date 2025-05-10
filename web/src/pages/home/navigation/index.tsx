import { A } from "@solidjs/router";
import CreateUpload from "@sv/components/create-upload";
import { createSignal, For, Show } from "solid-js";

const menuItems = [
  {
    label: "Home",
    name: "home",
    href: "/",
    icon: () => (
      <svg
        xmlns="http://www.w3.org/2000/svg"
        fill="none"
        viewBox="0 0 24 24"
        stroke-width="1.5"
        stroke="currentColor"
        class="size-6"
      >
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          d="m2.25 12 8.954-8.955c.44-.439 1.152-.439 1.591 0L21.75 12M4.5 9.75v10.125c0 .621.504 1.125 1.125 1.125H9.75v-4.875c0-.621.504-1.125 1.125-1.125h2.25c.621 0 1.125.504 1.125 1.125V21h4.125c.621 0 1.125-.504 1.125-1.125V9.75M8.25 21h8.25"
        />
      </svg>
    ),
  },
  {
    label: "Shared",
    name: "shared",
    href: "/shared",
    icon: () => (
      <svg
        xmlns="http://www.w3.org/2000/svg"
        fill="none"
        viewBox="0 0 24 24"
        stroke-width="1.5"
        stroke="currentColor"
        class="size-6"
      >
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          d="M18 18.72a9.094 9.094 0 0 0 3.741-.479 3 3 0 0 0-4.682-2.72m.94 3.198.001.031c0 .225-.012.447-.037.666A11.944 11.944 0 0 1 12 21c-2.17 0-4.207-.576-5.963-1.584A6.062 6.062 0 0 1 6 18.719m12 0a5.971 5.971 0 0 0-.941-3.197m0 0A5.995 5.995 0 0 0 12 12.75a5.995 5.995 0 0 0-5.058 2.772m0 0a3 3 0 0 0-4.681 2.72 8.986 8.986 0 0 0 3.74.477m.94-3.197a5.971 5.971 0 0 0-.94 3.197M15 6.75a3 3 0 1 1-6 0 3 3 0 0 1 6 0Zm6 3a2.25 2.25 0 1 1-4.5 0 2.25 2.25 0 0 1 4.5 0Zm-13.5 0a2.25 2.25 0 1 1-4.5 0 2.25 2.25 0 0 1 4.5 0Z"
        />
      </svg>
    ),
  },
  {
    label: "Favorites",
    name: "favorites",
    href: "/favorites",
    icon: () => (
      <svg
        xmlns="http://www.w3.org/2000/svg"
        fill="none"
        viewBox="0 0 24 24"
        stroke-width="1.5"
        stroke="currentColor"
        class="size-6"
      >
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          d="M11.48 3.499a.562.562 0 0 1 1.04 0l2.125 5.111a.563.563 0 0 0 .475.345l5.518.442c.499.04.701.663.321.988l-4.204 3.602a.563.563 0 0 0-.182.557l1.285 5.385a.562.562 0 0 1-.84.61l-4.725-2.885a.562.562 0 0 0-.586 0L6.982 20.54a.562.562 0 0 1-.84-.61l1.285-5.386a.562.562 0 0 0-.182-.557l-4.204-3.602a.562.562 0 0 1 .321-.988l5.518-.442a.563.563 0 0 0 .475-.345L11.48 3.5Z"
        />
      </svg>
    ),
  },
  {
    label: "Trash",
    name: "trash",
    href: "/trash",
    icon: () => (
      <svg
        xmlns="http://www.w3.org/2000/svg"
        fill="none"
        viewBox="0 0 24 24"
        stroke-width="1.5"
        stroke="currentColor"
        class="size-6"
      >
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          d="m14.74 9-.346 9m-4.788 0L9.26 9m9.968-3.21c.342.052.682.107 1.022.166m-1.022-.165L18.16 19.673a2.25 2.25 0 0 1-2.244 2.077H8.084a2.25 2.25 0 0 1-2.244-2.077L4.772 5.79m14.456 0a48.108 48.108 0 0 0-3.478-.397m-12 .562c.34-.059.68-.114 1.022-.165m0 0a48.11 48.11 0 0 1 3.478-.397m7.5 0v-.916c0-1.18-.91-2.164-2.09-2.201a51.964 51.964 0 0 0-3.32 0c-1.18.037-2.09 1.022-2.09 2.201v.916m7.5 0a48.667 48.667 0 0 0-7.5 0"
        />
      </svg>
    ),
  },
];

// Main navigation sidebar
function Navigation() {
  const [activeMenu, setActiveMenu] = createSignal<string>("");

  return (
    <>
      {/* Mobile bottom navigation */}
      <div class="md:hidden fixed bottom-0 inset-x-0 bg-white border-t border-border z-10">
        <nav class="flex items-center justify-around">
          <For each={menuItems}>
            {(item, index) => (
              <>
                <Show when={index() === 2}>
                  <div class="mx-2">
                    <CreateUpload />
                  </div>
                </Show>
                <A href={item.href} title={item.label}>
                  <div
                    class={`m-1 py-2 px-2 rounded-full link-no-underline ${
                      activeMenu() === item.name ? "bg-primary-light/30" : ""
                    }`}
                    onClick={() => setActiveMenu(item.name)}
                  >
                    {item.icon()}
                  </div>
                </A>
              </>
            )}
          </For>
        </nav>
      </div>

      {/* Desktop sidebar */}
      <aside class="hidden md:block w-64 h-screen bg-white border-r border-border fixed left-0 top-0 overflow-y-auto p-4">
        <h1 class="gradient-text text-center">SkyVault</h1>

        <CreateUpload />

        <nav class="space-y-1">
          <For each={menuItems}>
            {(item) => (
              <A href={item.href} title={item.label}>
                <div
                  class={`link-no-underline flex items-center gap-3 px-3 py-2 rounded-md font-medium ${
                    activeMenu() === item.name ? "bg-primary-light/30" : ""
                  }`}
                  onClick={() => setActiveMenu(item.name)}
                >
                  {item.icon()}
                  <span>{item.label}</span>
                </div>
              </A>
            )}
          </For>
        </nav>
      </aside>
    </>
  );
}

export default Navigation;
