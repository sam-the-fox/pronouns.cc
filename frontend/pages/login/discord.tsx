import { useEffect, useState } from "react";
import { useRouter } from "next/router";
import { useRecoilState } from "recoil";
import fetchAPI from "../../lib/fetch";
import { userState } from "../../lib/state";
import { APIError, MeUser, SignupResponse } from "../../lib/types";
import TextInput from "../../components/TextInput";
import Loading from "../../components/Loading";
import { stat } from "fs";
import Button, { ButtonStyle } from "../../components/Button";
import Notice from "../../components/Notice";

interface CallbackResponse {
  has_account: boolean;
  token?: string;
  user?: MeUser;

  discord?: string;
  ticket?: string;
  require_invite: boolean;
}

interface State {
  hasAccount: boolean;
  isLoading: boolean;
  token: string | null;
  user: MeUser | null;
  discord: string | null;
  ticket: string | null;
  error?: any;
  requireInvite: boolean;
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
    requireInvite: false,
  });
  const [formData, setFormData] = useState<{
    username: string;
    invite: string;
  }>({ username: "", invite: "" });

  useEffect(() => {
    if (!router.query.code || !router.query.state) {
      return;
    }
    if (state.ticket || state.token) {
      return;
    }

    fetchAPI<CallbackResponse>("/auth/discord/callback", "POST", {
      callback_domain: window.location.origin,
      code: router.query.code,
      state: router.query.state,
    })
      .then((resp) => {
        setState({
          hasAccount: resp.has_account,
          isLoading: false,
          token: resp.token || null,
          user: resp.user || null,
          discord: resp.discord || null,
          ticket: resp.ticket || null,
          requireInvite: resp.require_invite,
        });
      })
      .catch((e) => {
        setState({
          hasAccount: false,
          isLoading: false,
          error: e,
          token: null,
          user: null,
          discord: null,
          ticket: null,
          requireInvite: false,
        });
      });

    // we got a token + user, save it and return to the home page
    if (state.token) {
      window.localStorage.setItem("pronouns-token", state.token);
      setUser(state.user!);

      router.push("/");
    }
  }, [state.token, state.user, setState, router]);

  if (!state.ticket && !state.error) {
    return <Loading />;
  } else if (!state.ticket && state.error) {
    return (
      <Notice style={ButtonStyle.danger} header="Login error">
        <p>{state.error.message ?? state.error}</p>
        <p>Try again?</p>
      </Notice>
    );
  }

  // user needs to create an account
  const signup = async () => {
    try {
      const resp = await fetchAPI<SignupResponse>(
        "/auth/discord/signup",
        "POST",
        {
          ticket: state.ticket,
          username: formData.username,
          invite_code: formData.invite,
        }
      );

      setUser(resp.user);
      localStorage.setItem("pronouns-token", resp.token);

      router.push("/");
    } catch (e) {
      setState({ ...state, error: e });
    }
  };

  return (
    <>
      <h1 className="font-bold text-lg">Get started</h1>
      <p>
        You{"'"}ve logged in with Discord as{" "}
        <strong className="font-bold">{state.discord}</strong>.
      </p>

      {state.error && (
        <Notice style={ButtonStyle.danger} header="Create account error">
          <p>{state.error.message ?? state.error}</p>
          <p>Try again?</p>
        </Notice>
      )}

      <label>
        <span className="font-bold">Username</span>
        <TextInput
          contrastBackground
          value={formData.username}
          onChange={(e) =>
            setFormData({ ...formData, username: e.target.value })
          }
        />
      </label>
      {state.requireInvite && (
        <label>
          <span className="font-bold">Invite code</span>
          <TextInput
            contrastBackground
            value={formData.invite}
            onChange={(e) =>
              setFormData({ ...formData, invite: e.target.value })
            }
          />
        </label>
      )}
      <Button style={ButtonStyle.success} onClick={() => signup()}>
        Create account
      </Button>
    </>
  );
}
