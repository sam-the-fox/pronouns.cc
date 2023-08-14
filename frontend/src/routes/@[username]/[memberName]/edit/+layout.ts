import {
  type PrideFlag,
  type MeUser,
  type APIError,
  type Member,
  type PronounsJson,
  ErrorCode,
} from "$lib/api/entities";
import { apiFetchClient } from "$lib/api/fetch";
import { error, redirect } from "@sveltejs/kit";

import pronounsRaw from "$lib/pronouns.json";
import type { LayoutLoad } from "./$types";
const pronouns = pronounsRaw as PronounsJson;

export const ssr = false;

export const load = (async ({ params }) => {
  try {
    const user = await apiFetchClient<MeUser>(`/users/@me`);
    const member = await apiFetchClient<Member>(
      `/users/${params.username}/members/${params.memberName}`,
    );
    const flags = await apiFetchClient<PrideFlag[]>("/users/@me/flags");

    if (user.id !== member.user.id) {
      throw { code: ErrorCode.NotOwnMember, message: "Can only edit your own members" } as APIError;
    }

    if (
      user.name !== params.username ||
      member.user.name !== params.username ||
      member.name !== params.memberName
    ) {
      throw redirect(303, `/@${user.name}/${member.name}`);
    }

    return {
      user,
      member,
      pronouns: pronouns.autocomplete,
      flags,
    };
  } catch (e) {
    if ("code" in e) throw error(500, e as APIError);
    throw e;
  }
}) satisfies LayoutLoad;
