import { atom } from "recoil";
import { MeUser } from "./types";

export const userState = atom<MeUser | null>({
  key: "userState",
  default: null,
});
