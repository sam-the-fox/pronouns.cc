import type { PrideFlag, APIError, MeUser, PronounsJson } from "$lib/api/entities";
import { apiFetchClient } from "$lib/api/fetch";
import { error, redirect, type Redirect } from "@sveltejs/kit";

import pronounsRaw from "$lib/pronouns.json";
const pronouns = pronounsRaw as PronounsJson;

export const ssr = false;

export const load = async ({ params }) => {
  try {
    const user = await apiFetchClient<MeUser>(`/users/@me`);
    const flags = await apiFetchClient<PrideFlag[]>("/users/@me/flags");

    if (params.username !== user.name) {
      redirect(303, `/@${user.name}/edit`);
    }

    return {
      user,
      pronouns: pronouns.autocomplete,
      flags,
    };
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
  } catch (e: any) {
    if ("code" in e) error(500, e as App.Error);
    throw e;
  }
};
