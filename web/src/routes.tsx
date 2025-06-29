import { Navigate, Route, Router } from "@solidjs/router";
import AppLayout from "@sv/pages/appLayout";
import AuthLayout from "@sv/pages/auth/authLayout";
import SignIn from "@sv/pages/auth/signIn";
import SignUp from "@sv/pages/auth/signUp";
import Drive from "@sv/pages/drive";
import Favorites from "@sv/pages/favorites";
import Shared from "@sv/pages/shared";
import Trash from "@sv/pages/trash";
import { CLIENT_URLS } from "./utils/consts";

export default function AppRoutes() {
  return (
    <Router explicitLinks={true}>
      <Route path="/auth" component={AuthLayout}>
        <Route path="/sign-in" component={SignIn} />
        <Route path="/sign-up" component={SignUp} />
      </Route>
      <Route path="/" component={AppLayout}>
        <Route
          path=""
          component={() => <Navigate href={CLIENT_URLS.DRIVE} />}
        />
        <Route path={`${CLIENT_URLS.DRIVE}/:folderId?`} component={Drive} />
        <Route path={`${CLIENT_URLS.SHARED}/:folderId?`} component={Shared} />
        <Route path={CLIENT_URLS.FAVORITES} component={Favorites} />
        <Route path={CLIENT_URLS.TRASH} component={Trash} />
        <Route
          path="*"
          component={() => <Navigate href={CLIENT_URLS.DRIVE} />}
        />
      </Route>
    </Router>
  );
}
