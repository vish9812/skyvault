import authSvc from "@/services/auth.svc";
import { Button } from "./ui/button";
import { Search } from "lucide-react";
import FileUploader from "./FileUploader";

const Header = () => {
  return (
    <header className="header">
      <Search />
      <div className="header-wrapper">
        <FileUploader />
        <Button type="submit" className="sign-out-button" onClick={authSvc.signOut}>
            <img
              src="/assets/icons/logout.svg"
              alt="logo"
              width={24}
              height={24}
              className="w-6"
            />
          </Button>
      </div>
    </header>
  );
};
export default Header;
