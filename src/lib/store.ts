import { writable } from "svelte/store";
import { browser } from "$app/environment";

import type { MeUser } from "./api/entities";

export const userStore = writable<MeUser | null>(null);
export const tokenStore = writable<string | null>(null);

let defaultThemeValue = "dark";
const initialThemeValue = browser
  ? window.localStorage.getItem("pronouns-theme") ?? defaultThemeValue
  : defaultThemeValue;

export const themeStore = writable<string>(initialThemeValue);
