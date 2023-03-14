import type { APIError, Member } from "$lib/api/entities";
import { apiFetch } from "$lib/api/fetch";
import { error } from "@sveltejs/kit";

export const ssr = false;

export const load = async ({ params }) => {
  try {
    const member = await apiFetch<Member>(`/members/${params.id}`, {});

    return {
      member,
    };
  } catch (e) {
    throw error((e as APIError).code, (e as APIError).message);
  }
};
