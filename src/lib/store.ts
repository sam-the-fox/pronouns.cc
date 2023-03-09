import { writable } from "svelte/store";

import type { MeUser } from "./api/entities";

export const userStore = writable<MeUser | null>(null);
export const tokenStore = writable<string | null>(null);