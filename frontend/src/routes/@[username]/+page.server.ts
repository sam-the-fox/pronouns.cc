import { apiFetch } from "$lib/api/fetch";
import { ErrorCode, type APIError, type User } from "$lib/api/entities";
import { error } from "@sveltejs/kit";

export const load = async ({ params }) => {
  try {
    const resp = await apiFetch<User>(`/users/${params.username}`, {
      method: "GET",
    });

    return resp;
  } catch (e) {
    if ((e as APIError).code === ErrorCode.UserNotFound) {
      error(404, e as APIError);
    }

    throw e;
  }
};
