import { Html, Head, Main, NextScript } from "next/document";

export default function Document() {
  return (
    <Html>
      <Head>
        <link rel="icon" type="image/svg+xml" href="/favicon.svg" />
      </Head>
      <body className="bg-white dark:bg-slate-800 text-black dark:text-white">
        <Main />
        <NextScript />
      </body>
    </Html>
  );
}
