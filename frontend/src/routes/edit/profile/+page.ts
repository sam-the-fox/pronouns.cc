import { ErrorCode, type APIError, type MeUser } from "$lib/api/entities";
import { apiFetchClient } from "$lib/api/fetch";
import { error, redirect } from "@sveltejs/kit";

export const ssr = false;

export const load = async () => {
  try {
    const resp = await apiFetchClient<MeUser>(`/users/@me`);

    throw redirect(303, `/@${resp.name}/edit`);
  } catch (e) {
    if (
      (e as APIError).code === ErrorCode.Forbidden ||
      (e as APIError).code === ErrorCode.InvalidToken
    ) {
      throw error(403, e as APIError);
    }

    throw e;
  }
};
