import type { MeUser, APIError, Member } from "$lib/api/entities";
import { apiFetchClient } from "$lib/api/fetch";
import { error } from "@sveltejs/kit";

export const ssr = false;

export const load = async ({ params }) => {
  try {
    const user = await apiFetchClient<MeUser>(`/users/@me`);
    const member = await apiFetchClient<Member>(`/members/${params.id}`);

    return {
      user,
      member,
    };
  } catch (e) {
    throw error((e as APIError).code, (e as APIError).message);
  }
};
