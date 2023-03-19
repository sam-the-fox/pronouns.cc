import { apiFetch } from "$lib/api/fetch";
import { ErrorCode, type APIError, type Member } from "$lib/api/entities";
import { error } from "@sveltejs/kit";

export const load = async ({ params }) => {
  try {
    const resp = await apiFetch<Member>(`/users/${params.username}/members/${params.memberName}`, {
      method: "GET",
    });

    return resp;
  } catch (e) {
    if (
      (e as APIError).code === ErrorCode.UserNotFound ||
      (e as APIError).code === ErrorCode.MemberNotFound
    ) {
      throw error(404, (e as APIError).message);
    }

    throw error(500, (e as APIError).message);
  }
};
