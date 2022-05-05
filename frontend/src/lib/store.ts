import axios from "axios";
import { atom } from "recoil";
import type { APIError, MeUser } from "./types";

export const userState = atom<MeUser | null>({
  key: "userState",
  default: getCurrentUser(),
});

async function getCurrentUser(): Promise<MeUser | null> {
  const token = localStorage.getItem("pronouns-token");
  if (!token) return null;

  try {
    const resp = await axios.get<MeUser | APIError>("/api/v1/users/@me");
    if ("id" in resp.data) return resp.data as MeUser;
    return null;
  } catch (e) {
    console.log("Error fetching /users/@me:", e);
  }

  return null;
}
