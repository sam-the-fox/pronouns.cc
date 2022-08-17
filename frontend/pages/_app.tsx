import "../styles/globals.css";
import type { AppProps } from "next/app";
import Container from "../components/Container";
import Navigation from "../components/Navigation";
import { RecoilRoot } from "recoil";
import Head from "next/head";

function MyApp({ Component, pageProps }: AppProps) {
  return (
    <RecoilRoot>
      <Head>
        <title key="title">pronouns.cc</title>
      </Head>
      <Navigation />
      <Container>
        <Component {...pageProps} />
      </Container>
    </RecoilRoot>
  );
}

export default MyApp;
