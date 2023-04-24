import type { APIError, MeUser } from "$lib/api/entities";
import { apiFetch } from "$lib/api/fetch";
import type { PageServerLoad } from "./$types";
import { PUBLIC_BASE_URL } from "$env/static/public";

export const load = (async ({ url }) => {
  try {
    const resp = await apiFetch<CallbackResponse>("/auth/google/callback", {
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

  google?: string;
  ticket?: string;
  require_invite: boolean;
  require_captcha: boolean;

  is_deleted: boolean;
  deleted_at?: string;
  self_delete?: boolean;
  delete_reason?: string;
}
