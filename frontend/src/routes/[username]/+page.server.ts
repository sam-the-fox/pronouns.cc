import { ErrorCode, type APIError, type User } from "$lib/api/entities";
import { apiFetch } from "$lib/api/fetch";
import { error, redirect } from "@sveltejs/kit";

export const load = async ({ params }) => {
  try {
    const resp = await apiFetch<User>(`/users/${params.username}`, {
      method: "GET",
    });

    redirect(303, `/@${resp.name}`);
  } catch (e) {
    if ((e as APIError).code === ErrorCode.UserNotFound) {
      error(404, e as App.Error);
    }

    throw e;
  }
};
