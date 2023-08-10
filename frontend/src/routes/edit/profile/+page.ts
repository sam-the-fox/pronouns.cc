import type { MeUser } from "$lib/api/entities";
import { apiFetchClient } from "$lib/api/fetch";
import { redirect } from "@sveltejs/kit";

export const ssr = false;

export const load = async () => {
  const resp = await apiFetchClient<MeUser>(`/users/@me`);

  throw redirect(303, `/@${resp.name}/edit`);
};
