import Header from "@/components/Header";
import MobileNavigation from "@/components/MobileNavigation";
import Sidebar from "@/components/Sidebar";
import { Toaster } from "@/components/ui/toaster";
import consts from "@/lib/consts";
import authSvc from "@/services/auth.svc";
import { useEffect } from "react";
import { Outlet, useNavigate } from "react-router";

const HomeLayout = () => {
  const navigate = useNavigate();
  const profile = authSvc.profile();

  useEffect(() => {
    if (!profile) {
      navigate(consts.pageRoutes.signIn);
    }
  }, [navigate, profile]);

  if (!profile) {
    return null;
  }

  return (
    <main className="flex h-screen">
      <Sidebar {...profile} />

      <section className="flex h-full flex-1 flex-col">
        <MobileNavigation {...profile} />
        <Header />
        <div className="main-content">
          <Outlet />
        </div>
      </section>

      <Toaster />
    </main>
  );
};

export default HomeLayout;
