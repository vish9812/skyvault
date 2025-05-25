import { DropdownMenu } from "@kobalte/core/dropdown-menu";
import { getProfile, signOut } from "@sv/apis/auth";
import { useNavigate } from "@solidjs/router";
import { CLIENT_URLS } from "@sv/utils/consts";
import format from "@sv/utils/format";

function ProfileDropdown() {
  const navigate = useNavigate();
  const profile = getProfile()!;

  const handleLogout = () => {
    signOut();
    navigate(CLIENT_URLS.SIGN_IN, { replace: true });
  };

  return (
    <DropdownMenu>
      <DropdownMenu.Trigger
        class="w-10 h-10 rounded-full border border-border cursor-pointer"
        classList={{
          "btn btn-gradient btn-gradient-d-expanded": !profile.avatarBase64,
        }}
      >
        {profile.avatarBase64 ? (
          <img
            src={`data:image/png;base64,${profile.avatarBase64}`}
            alt="avatar"
          />
        ) : (
          <span class="font-bold text-lg">
            {format.initials(profile.fullName)}
          </span>
        )}
      </DropdownMenu.Trigger>
      <DropdownMenu.Portal>
        <DropdownMenu.Content class="bg-white rounded-lg shadow-md border border-border-strong min-w-[140px] mt-2">
          <DropdownMenu.Item class="dropdown-item" onSelect={handleLogout}>
            <span class="flex items-center gap-2">
              <svg
                xmlns="http://www.w3.org/2000/svg"
                fill="none"
                viewBox="0 0 24 24"
                stroke-width="1.5"
                stroke="currentColor"
                class="size-5 text-neutral-light"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  d="M8.25 9V5.25A2.25 2.25 0 0 1 10.5 3h6a2.25 2.25 0 0 1 2.25 2.25v13.5A2.25 2.25 0 0 1 16.5 21h-6a2.25 2.25 0 0 1-2.25-2.25V15m-3 0-3-3m0 0 3-3m-3 3H15"
                />
              </svg>
              Logout
            </span>
          </DropdownMenu.Item>
        </DropdownMenu.Content>
      </DropdownMenu.Portal>
    </DropdownMenu>
  );
}

export default ProfileDropdown;
