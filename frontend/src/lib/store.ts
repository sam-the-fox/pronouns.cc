import axios from "axios";
import { atom, useRecoilState, useRecoilValue } from "recoil";
import fetchAPI from "./fetch";
import { APIError, ErrorCode, MeUser } from "./types";

export const userState = atom<MeUser>({
  key: "userState",
  default: getCurrentUser(),
});

async function getCurrentUser() {
  const token = localStorage.getItem("pronouns-token");
  if (!token) return null;

  try {
    return await fetchAPI<MeUser>("/users/@me");
  } catch (e) {
    if (
      (e as APIError).code === ErrorCode.Forbidden ||
      (e as APIError).code === ErrorCode.InvalidToken
    ) {
      localStorage.removeItem("pronouns-token");
    }

    console.log("Error fetching /users/@me:", e);
  }

  return null;
}

export function isMeUser(id: string): boolean {
  const meUser = useRecoilValue(userState);
  return meUser && meUser.id === id;
}
