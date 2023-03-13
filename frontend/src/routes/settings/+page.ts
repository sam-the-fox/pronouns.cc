import {
  type Invite,
  type APIError,
  type MeUser,
  type PartialMember,
  ErrorCode,
} from "$lib/api/entities";
import { apiFetchClient } from "$lib/api/fetch";
import { error } from "@sveltejs/kit";

export const load = async () => {
  try {
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

    return { user, members, invites, invitesEnabled };
  } catch (e) {
    throw error(500, (e as APIError).message);
  }
};
