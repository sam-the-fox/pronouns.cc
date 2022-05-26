import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { useRecoilState } from "recoil";
import fetchAPI from "../../lib/fetch";
import Loading from "../../lib/Loading";
import { userState } from "../../lib/store";
import { MeUser } from "../../lib/types";

interface CallbackResponse {
  has_account: boolean;
  token?: string;
  user?: MeUser;

  discord?: string;
  ticket?: string;
  require_invite?: boolean;
}

export default function Discord() {
  const navigate = useNavigate();
  const params = new URLSearchParams(window.location.search);

  const [state, setState] = useState({
    hasAccount: false,
    isLoading: false,
    token: null,
    user: null,
    discord: null,
    ticket: null,
    error: null,
  });

  const [user, setUser] = useRecoilState(userState);

  useEffect(() => {
    if (state.isLoading) return;
    setState({ ...state, isLoading: true });

    fetchAPI<CallbackResponse>("/auth/discord/callback", "POST", {
      callback_domain: window.location.origin,
      code: params.get("code"),
      state: params.get("state"),
    }).then(
      (resp) => {
        setState({
          hasAccount: resp.has_account,
          isLoading: false,
          token: resp.token,
          user: resp.user,
          discord: resp.discord,
          ticket: resp.ticket,
          error: null,
        });

        console.log("token:", resp.token);
        localStorage.setItem("pronouns-token", resp.token);

        if (resp.user) setUser(resp.user as MeUser);
      },
      (err) => {
        console.log(err);
        setState({ ...state, error: err, isLoading: false });
      }
    );
  }, []);

  if (user) {
    // we got a token + user, save it and return to the home page
    navigate("/");
  }

  if (user || state.isLoading) {
    return <Loading />;
  }

  return <>wow such login</>;
}
