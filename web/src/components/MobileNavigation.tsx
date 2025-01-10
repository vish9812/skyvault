import {
  Sheet,
  SheetContent,
  SheetTitle,
  SheetTrigger,
} from "@/components/ui/sheet";
import consts from "@/lib/consts";
import { cn } from "@/lib/utils";
import { useState } from "react";
import { Link, useLocation } from "react-router";
import FileUploader from "./FileUploader";
import { Button } from "./ui/button";
import authSvc from "@/services/auth.svc";
import { Separator } from "./ui/separator";
import Avatar from "./Avatar";

interface Props {
  fullName: string;
  avatar: string;
  email: string;
}

const MobileNavigation = ({ fullName, avatar, email }: Props) => {
  const [open, setOpen] = useState(false);
  const { pathname } = useLocation();

  return (
    <header className="mobile-header">
      <h1>Mobile Navigation</h1>
      <img
        src="/assets/icons/logo-full-brand.svg"
        alt="logo"
        width={120}
        height={52}
        className="h-auto"
      />

      <Sheet open={open} onOpenChange={setOpen}>
        <SheetTrigger>
          <img
            src="/assets/icons/menu.svg"
            alt="Search"
            width={30}
            height={30}
          />
        </SheetTrigger>
        <SheetContent className="shad-sheet h-screen px-3">
          <SheetTitle>
            <div className="header-user">
              <Avatar avatar={avatar} className="header-user-avatar" />
              <div className="sm:hidden lg:block">
                <p className="subtitle-2 capitalize">{fullName}</p>
                <p className="caption">{email}</p>
              </div>
            </div>
            <Separator className="mb-4 bg-light-200/20" />
          </SheetTitle>

          <nav className="mobile-nav">
            <ul className="mobile-nav-list">
              {consts.navItems.map(({ url, name, icon }) => (
                <Link key={name} to={url} className="lg:w-full">
                  <li
                    className={cn(
                      "mobile-nav-item",
                      pathname === url && "shad-active"
                    )}
                  >
                    <img
                      src={icon}
                      alt={name}
                      width={24}
                      height={24}
                      className={cn(
                        "nav-icon",
                        pathname === url && "nav-icon-active"
                      )}
                    />
                    <p>{name}</p>
                  </li>
                </Link>
              ))}
            </ul>
          </nav>

          <Separator className="my-5 bg-light-200/20" />

          <div className="flex flex-col justify-between gap-5 pb-5">
            <FileUploader />
            <Button
              type="submit"
              className="mobile-sign-out-button"
              onClick={authSvc.signOut}
            >
              <img
                src="/assets/icons/logout.svg"
                alt="logo"
                width={24}
                height={24}
              />
              <p>Logout</p>
            </Button>
          </div>
        </SheetContent>
      </Sheet>
    </header>
  );
};

export default MobileNavigation;
