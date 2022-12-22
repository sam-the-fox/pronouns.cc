import { GetServerSideProps } from "next";
import Head from "next/head";
import { useRouter } from "next/router";
import { useRecoilValue } from "recoil";

import { fetchAPI } from "../../lib/api-fetch";
import { userState } from "../../lib/state";

interface URLsResponse {
  discord: string;
}

export default function Login({ urls }: { urls: URLsResponse }) {
  const router = useRouter();

  if (useRecoilValue(userState) !== null) {
    router.push("/");
  }

  return (
    <>
      <Head>
        <title key="title">Login - pronouns.cc</title>
      </Head>
      <a href={urls.discord}>Login with Discord</a>
    </>
  );
}

export const getServerSideProps: GetServerSideProps = async (context) => {
  try {
    const urls = await fetchAPI<URLsResponse>("/auth/urls", "POST", {
      callback_domain: process.env.DOMAIN,
    });

    return { props: { urls } };
  } catch (e) {
    console.log(e);

    return { notFound: true };
  }
};
