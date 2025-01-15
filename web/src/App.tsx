import { Routes, Route } from "react-router";
import "./App.css";
import AuthLayout from "./auth/layout";
import SignIn from "./auth/pages/sign-in";
import SignUp from "./auth/pages/sign-up";
import HomeLayout from "./home/layout";
import Media from "./home/pages/media";
import Home from "./home/pages/home";

function App() {
  return (
    <>
      <Routes>
        <Route path="/" element={<HomeLayout />}>
          <Route index element={<Home />} />
          <Route path="/media/:mediaType" element={<Media />} />
        </Route>
        <Route element={<AuthLayout />}>
          <Route index path="/sign-in" element={<SignIn />} />
          <Route path="/sign-up" element={<SignUp />} />
        </Route>
      </Routes>
    </>
  );
}

export default App;
