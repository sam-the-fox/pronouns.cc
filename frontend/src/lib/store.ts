import { writable } from "svelte/store";
import { browser } from "$app/environment";

import type { MeUser } from "./api/entities";

const initialUserValue = null;
export const userStore = writable<MeUser | null>(initialUserValue);

const defaultThemeValue = "dark";
const initialThemeValue = browser
  ? window.localStorage.getItem("pronouns-theme") ?? defaultThemeValue
  : defaultThemeValue;

export const themeStore = writable<string>(initialThemeValue);

export const CURRENT_CHANGELOG = "0.4.0";
