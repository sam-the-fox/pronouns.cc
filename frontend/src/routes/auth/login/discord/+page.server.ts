import type { APIError, MeUser } from "$lib/api/entities";
import { apiFetch } from "$lib/api/fetch";
import type { PageServerLoad, Actions } from "./$types";
import { PUBLIC_BASE_URL } from "$env/static/public";

export const load = (async ({ url }) => {
  try {
    const resp = await apiFetch<CallbackResponse>("/auth/discord/callback", {
      method: "POST",
      body: {
        callback_domain: PUBLIC_BASE_URL,
        code: url.searchParams.get("code"),
        state: url.searchParams.get("state"),
      },
    });

    return {
      ...resp,
    };
  } catch (e) {
    return { error: e as APIError };
  }
}) satisfies PageServerLoad;

interface CallbackResponse {
  has_account: boolean;
  token?: string;
  user?: MeUser;

  discord?: string;
  ticket?: string;
  require_invite: boolean;
}
