import type { MeUser, APIError, Member, PronounsJson } from "$lib/api/entities";
import { apiFetchClient } from "$lib/api/fetch";
import { error } from "@sveltejs/kit";

import pronounsRaw from "$lib/pronouns.json";
const pronouns = pronounsRaw as PronounsJson;

export const ssr = false;

export const load = async ({ params }) => {
  try {
    const user = await apiFetchClient<MeUser>(`/users/@me`);
    const member = await apiFetchClient<Member>(`/members/${params.id}`);

    return {
      user,
      member,
      pronouns: pronouns.autocomplete,
    };
  } catch (e) {
    throw error((e as APIError).code, (e as APIError).message);
  }
};
