import { Routes, Route } from "react-router-dom";
import "./App.css";
import Container from "./lib/Container";
import Navigation from "./lib/Navigation";
import Home from "./pages/Home";
import User from "./pages/User";

function App() {
  return (
    <>
      <Navigation />
      <Container>
        <Routes>
          <Route path="/" element={<Home />} />
          <Route path="/u/:username" element={<User />} />
        </Routes>
      </Container>
    </>
  );
}

export default App;
