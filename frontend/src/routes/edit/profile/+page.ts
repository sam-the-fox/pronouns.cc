import type { PrideFlag, APIError, MeUser, PronounsJson } from "$lib/api/entities";
import { apiFetchClient } from "$lib/api/fetch";
import { error } from "@sveltejs/kit";

import pronounsRaw from "$lib/pronouns.json";
const pronouns = pronounsRaw as PronounsJson;

export const ssr = false;

export const load = async () => {
  try {
    const user = await apiFetchClient<MeUser>(`/users/@me`);
    const flags = await apiFetchClient<PrideFlag[]>("/users/@me/flags");

    return {
      user,
      pronouns: pronouns.autocomplete,
      flags,
    };
  } catch (e) {
    throw error((e as APIError).code, (e as APIError).message);
  }
};
