import { atom } from "recoil";
import { MeUser } from "./api-fetch";

export const userState = atom<MeUser | null>({
  key: "userState",
  default: null,
});

export const themeState = atom<boolean>({
  key: "themeState",
  default: false,
});
