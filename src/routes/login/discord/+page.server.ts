import type { MeUser } from "$lib/api/entities";
import { apiFetch } from "$lib/api/fetch";
import type { PageServerLoad } from "./$types";
import { PUBLIC_BASE_URL } from "$env/static/public";

export const load = (async (event) => {
  try {
    const resp = await apiFetch<CallbackResponse>("/auth/discord/callback", {
      method: "POST",
      body: {
        callback_domain: PUBLIC_BASE_URL,
        code: event.url.searchParams.get("code"),
        state: event.url.searchParams.get("state"),
      },
    });

    return {
      ...resp,
    };
  } catch (e) {
    return { error: e };
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
