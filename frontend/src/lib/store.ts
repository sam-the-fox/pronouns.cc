import { writable } from "svelte/store";
import { browser } from "$app/environment";

import type { MeUser } from "./api/entities";

const initialUserValue = null;
export const userStore = writable<MeUser | null>(initialUserValue);

let defaultThemeValue = "dark";
const initialThemeValue = browser
  ? window.localStorage.getItem("pronouns-theme") ?? defaultThemeValue
  : defaultThemeValue;

export const themeStore = writable<string>(initialThemeValue);
