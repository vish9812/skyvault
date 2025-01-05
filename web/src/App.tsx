import { Routes, Route } from "react-router";
import "./App.css";
import Layout from "./auth/layout";
import SignIn from "./auth/pages/sign-in";
import SignUp from "./auth/pages/sign-up";

function App() {
  return (
    <>
      {/* Render Auth Layout and Auth Routes */}
      <Routes>
        <Route element={<Layout />}>
          <Route index path="/sign-in" element={<SignIn />} />
          <Route path="/sign-up" element={<SignUp />} />
        </Route>
      </Routes>
    </>
  );
}

export default App;
