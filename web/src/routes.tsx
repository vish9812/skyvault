import { Route, Router } from "@solidjs/router";
import AuthLayout from "@sv/pages/auth/layout";
import SignUp from "@sv/pages/auth/signup";
import SignIn from "@sv/pages/auth/signin";
import AppLayout from "@sv/pages/app/layout";
import Home from "@sv/pages/app/home";

export default function AppRoutes() {
  return (
    <Router>
      <Route path="/auth" component={AuthLayout}>
        <Route path="/sign-in" component={SignIn} />
        <Route path="/sign-up" component={SignUp} />
      </Route>
      <Route path="/" component={AppLayout}>
        <Route path="/" component={Home} />
      </Route>
    </Router>
  );
}
