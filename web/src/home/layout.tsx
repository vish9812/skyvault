import consts from "@/lib/consts";
import authSvc from "@/services/auth.svc";
import React, { useEffect } from "react";
import { useNavigate } from "react-router";

const HomeLayout = () => {
  // If already authenticated, redirect to home
  const navigate = useNavigate();

  useEffect(() => {
    if (!authSvc.isAuthenticated()) {
      navigate(consts.pageRoutes.signIn);
    }
  }, [navigate]);

  if (!authSvc.isAuthenticated()) {
    return null;
  }

  return <div>Home Layout</div>;
};

export default HomeLayout;
