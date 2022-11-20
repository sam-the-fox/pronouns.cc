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
        <meta property="og:type" content="website" />
        <meta name="theme-color" content="#aa8ed6" />

        <meta key="sitename" property="og:site_name" content="pronouns.cc" />
        <meta key="title" property="og:title" content="pronouns.cc" />
        <meta
          key="description"
          property="og:description"
          content="Name and pronoun cards!"
        />
        <meta key="url" property="og:url" content={process.env.DOMAIN} />
      </Head>
      <Navigation />
      <Container>
        <Component {...pageProps} />
      </Container>
    </RecoilRoot>
  );
}

export default MyApp;
