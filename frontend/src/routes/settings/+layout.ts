import { ErrorCode, type APIError, type Invite, type MeUser } from "$lib/api/entities";
import { apiFetchClient } from "$lib/api/fetch";
import type { LayoutLoad } from "./$types";

export const ssr = false;

export const load = (async ({ parent }) => {
  const user = await apiFetchClient<MeUser>("/users/@me");

  let invites: Invite[] = [];
  let invitesEnabled = true;
  try {
    invites = await apiFetchClient<Invite[]>("/auth/invites");
  } catch (e) {
    if ((e as APIError).code === ErrorCode.InvitesDisabled) {
      invitesEnabled = false;
    }
  }

  const data = await parent();

  return {
    ...data,
    user,
    invites,
    invitesEnabled,
  };
}) satisfies LayoutLoad;
