import { Routes, Route } from "react-router";
import "./App.css";
import AuthLayout from "./auth/layout";
import SignIn from "./auth/pages/sign-in";
import SignUp from "./auth/pages/sign-up";
import HomeLayout from "./home/layout";

function App() {
  return (
    <>
      <Routes>
        <Route index element={<HomeLayout />} />
        <Route element={<AuthLayout />}>
          <Route index path="/sign-in" element={<SignIn />} />
          <Route path="/sign-up" element={<SignUp />} />
        </Route>
      </Routes>
    </>
  );
}

export default App;
