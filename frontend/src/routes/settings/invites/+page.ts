import { ErrorCode, type APIError, type Invite } from "$lib/api/entities";
import { apiFetchClient } from "$lib/api/fetch";
import { error } from "@sveltejs/kit";
import type { PageLoad } from "../$types";

export const load = (async () => {
  const data = {
    invitesEnabled: true,
    invites: [] as Invite[],
  };

  try {
    const invites = await apiFetchClient<Invite[]>("/auth/invites");
    data.invites = invites;
  } catch (e) {
    if ((e as APIError).code === ErrorCode.InvitesDisabled) {
      data.invitesEnabled = false;
      data.invites = [];
    } else {
      throw error((e as APIError).code, (e as APIError).message);
    }
  }

  return data;
}) satisfies PageLoad;
