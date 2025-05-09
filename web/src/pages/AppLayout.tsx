import { onMount, ParentComponent } from "solid-js";
import { Button } from "@kobalte/core/button";
import { DropdownMenu } from "@kobalte/core/dropdown-menu";
import { Toast } from "@kobalte/core/toast";
import { Portal } from "solid-js/web";
import { useNavigate } from "@solidjs/router";

// Toast notifications region
function Notifications() {
  return (
    <Portal>
      <Toast.Region class="fixed top-4 right-4 z-[100]">
        <Toast.List>
          <Toast
            toastId={1}
            class="bg-white shadow-lg rounded-lg overflow-hidden w-72"
          >
            <div class="p-4 flex items-start gap-3">
              <div class="flex-1">
                <Toast.Title class="font-medium text-gray-900">
                  Welcome to SkyVault!
                </Toast.Title>
                <Toast.Description class="text-sm text-gray-600 mt-1">
                  Your secure cloud storage solution
                </Toast.Description>
              </div>
              <Toast.CloseButton class="text-gray-400 hover:text-gray-600">
                <span class="material-symbols-outlined text-base">close</span>
              </Toast.CloseButton>
            </div>
            <Toast.ProgressTrack class="h-1 bg-gray-100">
              <Toast.ProgressFill class="bg-blue-500 h-1" />
            </Toast.ProgressTrack>
          </Toast>
        </Toast.List>
      </Toast.Region>
    </Portal>
  );
}

// User profile menu
function UserProfileMenu() {
  return (
    <DropdownMenu>
      <DropdownMenu.Trigger>
        <Button class="flex items-center gap-2 px-2 py-1 rounded-full hover:bg-gray-100">
          <span class="material-symbols-outlined text-gray-600">
            account_circle
          </span>
          <span class="text-gray-700 text-sm hidden sm:inline">User</span>
        </Button>
      </DropdownMenu.Trigger>
      <DropdownMenu.Content class="bg-white rounded-lg shadow-lg border border-gray-200 py-1 min-w-[180px] mt-1">
        <DropdownMenu.Item class="px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 cursor-pointer flex items-center gap-2">
          <span class="material-symbols-outlined text-gray-500 text-sm">
            settings
          </span>
          <span>Settings</span>
        </DropdownMenu.Item>
        <DropdownMenu.Item class="px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 cursor-pointer flex items-center gap-2">
          <span class="material-symbols-outlined text-gray-500 text-sm">
            logout
          </span>
          <span>Logout</span>
        </DropdownMenu.Item>
      </DropdownMenu.Content>
    </DropdownMenu>
  );
}

// Main navigation sidebar
function Navigation() {
  return (
    <>
      {/* Desktop sidebar */}
      <aside class="hidden md:block w-64 h-screen bg-white border-r border-gray-200 fixed left-0 top-0 z-20 overflow-y-auto">
        <div class="p-4">
          <div class="text-2xl font-extrabold bg-primary bg-clip-text text-transparent mb-8 select-none">
            SkyVault
          </div>
          <nav class="space-y-1">
            <a
              href="/"
              class="flex items-center gap-3 px-3 py-2 text-gray-800 rounded-md bg-primary/20 font-medium"
            >
              <span class="material-symbols-outlined text-primary">home</span>
              <span>Home</span>
            </a>
            <a
              href="#recents"
              class="flex items-center gap-3 px-3 py-2 text-gray-600 rounded-md hover:bg-primary/10"
            >
              <span class="material-symbols-outlined">history</span>
              <span>Recents</span>
            </a>
            <a
              href="#shared"
              class="flex items-center gap-3 px-3 py-2 text-gray-600 rounded-md hover:bg-primary/10"
            >
              <span class="material-symbols-outlined">share</span>
              <span>Shared</span>
            </a>
            <a
              href="#favorites"
              class="flex items-center gap-3 px-3 py-2 text-gray-600 rounded-md hover:bg-primary/10"
            >
              <span class="material-symbols-outlined">star</span>
              <span>Favorites</span>
            </a>
            <a
              href="#trash"
              class="flex items-center gap-3 px-3 py-2 text-gray-600 rounded-md hover:bg-primary/10"
            >
              <span class="material-symbols-outlined">delete</span>
              <span>Trash</span>
            </a>
          </nav>
        </div>
      </aside>

      {/* Mobile bottom navigation */}
      <nav class="fixed bottom-0 inset-x-0 bg-white border-t border-gray-200 flex justify-around md:hidden items-center z-50">
        <a href="/" class="flex flex-col items-center py-2 text-primary">
          <span class="material-symbols-outlined">home</span>
          <span class="text-xs">Home</span>
        </a>
        <a
          href="#recents"
          class="flex flex-col items-center py-2 text-gray-600 hover:text-primary"
        >
          <span class="material-symbols-outlined">history</span>
          <span class="text-xs">Recents</span>
        </a>
        <a
          href="#shared"
          class="flex flex-col items-center py-2 text-gray-600 hover:text-primary"
        >
          <span class="material-symbols-outlined">share</span>
          <span class="text-xs">Shared</span>
        </a>
        <a
          href="#favorites"
          class="flex flex-col items-center py-2 text-gray-600 hover:text-primary"
        >
          <span class="material-symbols-outlined">star</span>
          <span class="text-xs">Favorites</span>
        </a>
      </nav>
    </>
  );
}

const AppLayout: ParentComponent = (props) => {
  const navigate = useNavigate();
  onMount(() => {
    const token = localStorage.getItem("token");
    if (!token) {
      navigate("/sign-in", { replace: true });
    }
  });

  return (
    <div class="bg-gray-50 min-h-screen">
      {/* Navigation */}
      <Navigation />

      {/* Main Content */}
      <div class="md:pl-64 pt-16 pb-16 min-h-screen">
        {/* Header */}
        <header class="fixed top-0 left-0 right-0 md:left-64 bg-white h-16 border-b border-gray-200 z-30 px-4 flex items-center justify-between">
          <div class="md:hidden text-xl font-extrabold bg-primary bg-clip-text text-transparent select-none">
            SkyVault
          </div>
          <div class="flex items-center gap-4">
            <button class="p-2 rounded-full hover:bg-gray-100">
              <span class="material-symbols-outlined text-gray-600">
                search
              </span>
            </button>
            <UserProfileMenu />
          </div>
        </header>

        {/* Page Content */}
        <main class="p-4 md:p-6">{props.children}</main>
      </div>

      {/* Notifications */}
      <Notifications />
    </div>
  );
};

export default AppLayout;
