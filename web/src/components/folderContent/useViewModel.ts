import { useNavigate } from "@solidjs/router";
import { CLIENT_URLS } from "@sv/utils/consts";

function useViewModel() {
  const navigate = useNavigate();

  const handleFolderNavigation = (id: number) => {
    navigate(`${CLIENT_URLS.HOME}${id}`);
  };

  return {
    handleFolderNavigation,
  };
}

export default useViewModel;
