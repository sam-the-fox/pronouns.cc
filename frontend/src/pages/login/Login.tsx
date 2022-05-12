import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { useRecoilValue } from "recoil";
import fetchAPI from "../../lib/fetch";
import { userState } from "../../lib/store";

interface URLsResponse {
  discord: string;
}

export default function Login() {
  const [state, setState] = useState({
    loading: false,
    error: null,
    discord: "",
  });

  if (useRecoilValue(userState) !== null) {
    const nav = useNavigate();
    nav("/");
  }

  useEffect(() => {
    if (state.loading) return;
    setState({ ...state, loading: true });

    fetchAPI<URLsResponse>("/auth/urls", "POST", {
      callback_domain: window.location.origin,
    }).then(
      (resp) => {
        setState({ loading: false, error: null, discord: resp.discord });
      },
      (err) => {
        console.log(err);
        setState({ ...state, loading: false, error: err });
      }
    );
  }, []);

  if (state.loading) {
    return <>Loading...</>;
  } else if (state.error) {
    return <>Error: {`${state.error}`}</>;
  }

  return (
    <>
      <a href={state.discord}>Login with Discord</a>
    </>
  );
}
