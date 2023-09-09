import { writable } from "svelte/store";
import { browser } from "$app/environment";

import type { MeUser, Settings } from "./api/entities";

const initialUserValue = null;
export const userStore = writable<MeUser | null>(initialUserValue);

const defaultThemeValue = "dark";
const initialThemeValue = browser
  ? window.localStorage.getItem("pronouns-theme") ?? defaultThemeValue
  : defaultThemeValue;

export const themeStore = writable<string>(initialThemeValue);

const defaultSettingsValue = {
  settings: { read_changelog: "0.0.0", read_settings_notice: "0" } as Settings,
  current: false,
};
export const settingsStore = writable(defaultSettingsValue);

export const CURRENT_CHANGELOG = "0.6.0";
