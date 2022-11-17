import { useEffect, useState } from "react";
import { useRouter } from "next/router";
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
  token: string | null;
  user: MeUser | null;
  discord: string | null;
  ticket: string | null;
  error?: any;
}

export default function Discord() {
  const router = useRouter();

  const [user, setUser] = useRecoilState(userState);
  const [state, setState] = useState<State>({
    hasAccount: false,
    isLoading: false,
    token: null,
    user: null,
    discord: null,
    ticket: null,
    error: null,
  });

  useEffect(() => {
    if (!router.query.code || !router.query.state) { return; }

    fetchAPI<CallbackResponse>(
      "/auth/discord/callback",
      "POST",
      {
        callback_domain: window.location.origin,
        code: router.query.code,
        state: router.query.state,
      }
    ).then(resp => {
      setState({
        hasAccount: resp.has_account,
        isLoading: false,
        token: resp.token || null,
        user: resp.user || null,
        discord: resp.discord || null,
        ticket: resp.ticket || null,
      })
    }).catch(e => {
      return {
        props: {
          hasAccount: false,
          isLoading: false,
          error: e,
          token: null,
          user: null,
          discord: null,
          ticket: null,
        },
      };
    })

    // we got a token + user, save it and return to the home page
    if (state.token) {
      window.localStorage.setItem("pronouns-token", state.token);
      setUser(state.user!);

      router.push("/");
    }
  }, [state.token, state.user, setState, router]);

  return <>wow such login</>;
}
