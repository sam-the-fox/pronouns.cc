import axios from "axios";
import { atom, useRecoilState, useRecoilValue } from "recoil";
import { APIError, ErrorCode, MeUser } from "./types";

export const userState = atom<MeUser>({
  key: "userState",
  default: getCurrentUser(),
});

async function getCurrentUser() {
  const token = localStorage.getItem("pronouns-token");
  if (!token) return null;

  try {
    const resp = await axios.get<MeUser | APIError>("/api/v1/users/@me");
    if (resp.status === 200) {
      return resp.data as MeUser;
    }

    // if we got a forbidden error, the token is invalid
    if ((resp.data as APIError).code === ErrorCode.Forbidden) {
      localStorage.removeItem("pronouns-token");
    }
  } catch (e) {
    console.log("Error fetching /users/@me:", e);
  }

  return null;
}

export function isMeUser(id: string): boolean {
  const meUser = useRecoilValue(userState);
  return meUser && meUser.id === id;
}
