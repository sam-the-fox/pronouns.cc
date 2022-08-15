import "../styles/globals.css";
import type { AppProps } from "next/app";
import Container from "../components/Container";
import Navigation from "../components/Navigation";
import { RecoilRoot } from "recoil";

function MyApp({ Component, pageProps }: AppProps) {
  return (
    <RecoilRoot>
      <Navigation />
      <Container>
        <Component {...pageProps} />
      </Container>
    </RecoilRoot>
  );
}

export default MyApp;
