import { Routes, Route, useParams } from "react-router-dom";
import "./App.css";
import Container from "./lib/Container";
import Navigation from "./lib/Navigation";
import EditMe from "./pages/EditMe";
import Home from "./pages/Home";
import Discord from "./pages/login/Discord";
import Login from "./pages/login/Login";
import User from "./pages/User";

function App() {
  return (
    <>
      <Navigation />
      <Container>
        <Routes>
          <Route path="/" element={<Home />} />
          <Route path="/@:username" element={<User />} />
          <Route path="/@:username/:member" element={<User />} />
          <Route path="/edit" element={<EditMe />} />
          <Route path="/edit/:member" element={<EditMe />} />
          <Route path="/login" element={<Login />} />
          <Route path="/login/discord" element={<Discord />} />
        </Routes>
      </Container>
    </>
  );
}

export default App;
