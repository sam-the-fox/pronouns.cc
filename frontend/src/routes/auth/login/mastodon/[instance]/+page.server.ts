import type { APIError, MeUser } from "$lib/api/entities";
import { apiFetch } from "$lib/api/fetch";
import type { PageServerLoad } from "./$types";

export const load = (async ({ url, params }) => {
  try {
    const resp = await apiFetch<CallbackResponse>("/auth/mastodon/callback", {
      method: "POST",
      body: {
        instance: params.instance,
        code: url.searchParams.get("code"),
        state: url.searchParams.get("state"),
      },
    });

    return {
      ...resp,
      instance: params.instance,
    };
  } catch (e) {
    return { error: e as APIError };
  }
}) satisfies PageServerLoad;

interface CallbackResponse {
  has_account: boolean;
  token?: string;
  user?: MeUser;

  fediverse?: string;
  ticket?: string;
  require_invite: boolean;

  is_deleted: boolean;
  deleted_at?: string;
  self_delete?: boolean;
  delete_reason?: string;
}
