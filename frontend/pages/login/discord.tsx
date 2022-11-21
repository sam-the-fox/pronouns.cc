import { useEffect, useState } from "react";
import { useRouter } from "next/router";
import { useRecoilState } from "recoil";
import fetchAPI from "../../lib/fetch";
import { userState } from "../../lib/state";
import { APIError, MeUser, SignupResponse } from "../../lib/types";
import TextInput from "../../components/TextInput";
import Loading from "../../components/Loading";
import Button, { ButtonStyle } from "../../components/Button";
import Notice from "../../components/Notice";
import BlueLink from "../../components/BlueLink";
import toast from "../../lib/toast";

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
    if (state.isLoading || !router.query.code || !router.query.state) {
      return;
    }

    // we got a token + user, save it and return to the home page
    if (state.token) {
      window.localStorage.setItem("pronouns-token", state.token);
      setUser(state.user!);

      router.push("/");
      return;
    }

    setState({ ...state, isLoading: true });
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

    if (!state.ticket && !state.token) {
      return;
    }
  }, [router.query.code, router.query.state, state.token]);

  if (state.isLoading || (!state.ticket && !state.error)) {
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

      toast({ text: "Created account!", background: "success" });
      router.push("/");
    } catch (e: any) {
      toast({ text: `${e.message ?? e}`, background: "error" });
    }
  };

  return (
    <div>
      <div className="border-slate-200 dark:border-slate-700 border rounded max-w-xl">
        <div className="border-b border-slate-200 dark:border-slate-700 p-2">
          <h1 className="font-bold text-xl">Get started</h1>
        </div>
        <div className="px-2 pt-2">
          <p>
            Just one more thing!
            <br />
            <strong className="font-bold">{state.discord}</strong>, you need to
            choose a username
            {state.requireInvite && " and fill in an invite code"} before you
            create an account!
          </p>
        </div>
        <label className="block px-2 pb-3 pt-2">
          <span className="block font-bold p-2 text-slate-800 dark:text-slate-200">
            Username
          </span>
          <TextInput
            contrastBackground
            value={formData.username}
            onChange={(e) =>
              setFormData({ ...formData, username: e.target.value })
            }
          />
        </label>
        {state.requireInvite && (
          <label className="block px-2 pb-3 pt-2">
            <span className="block p-2 text-slate-800 dark:text-slate-200">
              Invite code <span className="font-bold">Invite code</span>
              <span className="font-italic">
                (an existing user can give you this!)
              </span>
            </span>
            <TextInput
              contrastBackground
              value={formData.invite}
              onChange={(e) =>
                setFormData({ ...formData, invite: e.target.value })
              }
            />
          </label>
        )}
        <div className="bg-slate-100 dark:bg-slate-700 border-t border-slate-200 dark:border-slate-700">
          <span className="block p-3">
            <Button style={ButtonStyle.success} onClick={() => signup()}>
              Create account
            </Button>
          </span>
          <span className="block px-3 pb-3">
            <span className="font-bold">Note:</span> by clicking &quot;Create
            account&quot;, you agree to the{" "}
            <BlueLink to="/page/tos">terms of service</BlueLink> and the{" "}
            <BlueLink to="/page/privacy">privacy policy</BlueLink>.
          </span>
        </div>
      </div>
    </div>
  );
}
