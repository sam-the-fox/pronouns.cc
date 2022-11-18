import { useRouter } from "next/router";
import { useEffect } from "react";
import { useRecoilState } from "recoil";
import Loading from "../components/Loading";
import { userState } from "../lib/state";

export default function Logout() {
  const router = useRouter();
  const [_, setUser] = useRecoilState(userState);

  useEffect(() => {
    localStorage.removeItem("pronouns-token");
    setUser(null);

    router.push("/");
  }, []);

  return <Loading />;
}
