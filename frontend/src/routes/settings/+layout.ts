import {
  ErrorCode,
  type Warning,
  type APIError,
  type Invite,
  type MeUser,
  type Member,
} from "$lib/api/entities";
import { apiFetchClient } from "$lib/api/fetch";
import { redirect } from "@sveltejs/kit";
import type { LayoutLoad } from "./$types";

export const ssr = false;

export const load = (async ({ parent }) => {
  try {
    const user = await apiFetchClient<MeUser>("/users/@me");
    const members = await apiFetchClient<Member[]>("/users/@me/members");
    const warnings = await apiFetchClient<Warning[]>("/auth/warnings?all=true");

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
      warnings,
    };
  } catch (e) {
    if ((e as APIError).code !== ErrorCode.InternalServerError) {
      redirect(303, "/auth/login");
    }
    throw e;
  }
}) satisfies LayoutLoad;
