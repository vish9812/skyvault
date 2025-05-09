import CreateUpload from "@sv/components/create-upload";

// Main navigation sidebar
function Navigation() {
  return (
    <>
      {/* Desktop sidebar */}
      <aside class="hidden md:block w-64 h-screen bg-bg border-r border-border fixed left-0 top-0 overflow-y-auto">
        <div class="p-4">
          <h2 class="gradient-text text-center">SkyVault</h2>

          <CreateUpload />

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

export default Navigation;
