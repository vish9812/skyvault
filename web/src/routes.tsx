import { Route, Router } from "@solidjs/router";
import SignIn from "@sv/pages/SignIn";
import SignUp from "@sv/pages/SignUp";
import Home from "@sv/pages/Home";
import AuthLayout from "@sv/pages/AuthLayout";

export default function AppRoutes() {
  return (
    <Router>
      <Route
        path="/sign-in"
        component={() => (
          <AuthLayout title="Sign In">
            <SignIn />
          </AuthLayout>
        )}
      />
      <Route
        path="/sign-up"
        component={() => (
          <AuthLayout title="Sign Up">
            <SignUp />
          </AuthLayout>
        )}
      />
      <Route path="/" component={Home} />
    </Router>
  );
}
