import { Route, Router } from "@solidjs/router";
// import SignIn from "@sv/pages/auth/signin";
import SignUp from "@sv/pages/auth/signup";
import Home from "@sv/pages/Home";
import AuthLayout from "@sv/pages/auth/layout";
import AppLayout from "./pages/AppLayout";

export default function AppRoutes() {
  return (
    <Router>
      <Route path="/auth" component={AuthLayout}>
        {/* <Route path="/sign-in" component={SignIn} /> */}
        <Route path="/sign-up" component={SignUp} />
      </Route>
      <Route path="/" component={AppLayout}>
        <Route path="/" component={Home} />
      </Route>
    </Router>
  );
}
