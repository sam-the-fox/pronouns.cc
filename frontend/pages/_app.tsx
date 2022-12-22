import "../styles/globals.css";

import type { AppProps } from "next/app";
import Head from "next/head";
import { RecoilRoot } from "recoil";

import Container from "../components/Container";
import Navigation from "../components/Navigation";

function MyApp({ Component, pageProps }: AppProps) {
  const domain =
    typeof window !== "undefined" ? window.location.origin : process.env.DOMAIN;

  return (
    <RecoilRoot>
      <Head>
        <title key="title">pronouns.cc</title>
        <meta property="og:type" content="website" />
        <meta name="theme-color" content="#aa8ed6" />

        <meta key="og:sitename" property="og:site_name" content="pronouns.cc" />
        <meta key="og:title" property="og:title" content="pronouns.cc" />
        <meta
          key="og:description"
          property="og:description"
          content="Name and pronoun cards!"
        />
        <meta key="og:url" property="og:url" content={domain} />
      </Head>
      <Navigation />
      <Container>
        <Component {...pageProps} />
      </Container>
    </RecoilRoot>
  );
}

export default MyApp;
