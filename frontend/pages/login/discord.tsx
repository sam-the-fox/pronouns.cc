import { useEffect } from "react";
import { useRouter } from "next/router";
import { GetServerSideProps } from "next";
import { useRecoilState } from "recoil";
import fetchAPI from "../../lib/fetch";
import { userState } from "../../lib/state";
import { MeUser } from "../../lib/types";

interface CallbackResponse {
  has_account: boolean;
  token?: string;
  user?: MeUser;

  discord?: string;
  ticket?: string;
  require_invite?: boolean;
}

interface State {
  hasAccount: boolean;
  isLoading: boolean;
  token?: string;
  user?: MeUser;
  discord?: string;
  ticket?: string;
  error?: any;
}

export default function Discord(props: State) {
  const router = useRouter();

  const [user, setUser] = useRecoilState(userState);

  useEffect(() => {
    // we got a token + user, save it and return to the home page
    if (props.token) {
      localStorage.setItem("pronouns-token", props.token);
      setUser(props.user!);

      router.push("/");
    }
  }, [props.token, props.user, setUser, router]);

  return <>wow such login</>;
}

export const getServerSideProps: GetServerSideProps<State> = async (
  context
) => {
  try {
    const resp = await fetchAPI<CallbackResponse>(
      "/auth/discord/callback",
      "POST",
      {
        callback_domain: process.env.DOMAIN,
        code: context.query.code,
        state: context.query.state,
      }
    );

    return {
      props: {
        hasAccount: resp.has_account,
        isLoading: false,
        token: resp.token,
        user: resp.user,
        discord: resp.discord || null,
        ticket: resp.ticket || null,
      },
    };
  } catch (e) {
    return { props: { error: e } };
  }
};
