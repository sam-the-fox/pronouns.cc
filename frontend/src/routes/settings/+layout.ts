import {
  ErrorCode,
  type APIError,
  type Invite,
  type MeUser,
  type PartialMember,
} from "$lib/api/entities";
import { apiFetchClient } from "$lib/api/fetch";
import type { LayoutLoad } from "./$types";

export const ssr = false;

export const load = (async ({ parent }) => {
  const user = await apiFetchClient<MeUser>("/users/@me");
  const members = await apiFetchClient<PartialMember[]>("/users/@me/members");

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
    members,
    invites,
    invitesEnabled,
  };
}) satisfies LayoutLoad;
