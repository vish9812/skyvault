import ProfileDropdown from "@sv/components/profileDropdown";
import Search from "@sv/components/search";

function AppHeader() {
  return (
    <header class="fixed top-0 left-0 right-0 md:ml-64 bg-white border-b border-border p-2">
      <div class="flex items-center justify-between">
        <h2 class="gradient-text md:hidden">SkyVault</h2>
        <div>
          <Search />
        </div>
        <ProfileDropdown />
      </div>
    </header>
  );
}

export default AppHeader;
