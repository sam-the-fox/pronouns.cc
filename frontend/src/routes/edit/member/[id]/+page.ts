import { ErrorCode, type APIError, type Member, type MeUser } from "$lib/api/entities";
import { apiFetchClient } from "$lib/api/fetch";
import { error, redirect } from "@sveltejs/kit";

export const ssr = false;

export const load = async ({ params }) => {
  try {
    const user = await apiFetchClient<MeUser>(`/users/@me`);
    const member = await apiFetchClient<Member>(`/members/${params.id}`);

    if (user.id !== member.user.id) {
      throw {
        code: ErrorCode.NotOwnMember,
        message: "You can only edit your own members.",
      } as APIError;
    }

    throw redirect(303, `/@${user.name}/${member.name}/edit`);
  } catch (e) {
    if (
      (e as APIError).code === ErrorCode.Forbidden ||
      (e as APIError).code === ErrorCode.InvalidToken ||
      (e as APIError).code === ErrorCode.NotOwnMember
    ) {
      throw error(403, e as APIError);
    }

    throw e;
  }
};
