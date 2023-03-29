import type { APIError, MeUser } from "$lib/api/entities";
import { apiFetchClient } from "$lib/api/fetch";
import { error } from "@sveltejs/kit";

export const ssr = false;

export const load = async () => {
  try {
    const user = await apiFetchClient<MeUser>(`/users/@me`);

    return {
      user,
    };
  } catch (e) {
    throw error((e as APIError).code, (e as APIError).message);
  }
};
