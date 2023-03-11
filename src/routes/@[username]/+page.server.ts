import { apiFetch } from "$lib/api/fetch";
import type { User } from "$lib/api/entities";

export const load = async ({ params }) => {
  const resp = await apiFetch<User>(`/users/${params.username}`, {
    method: "GET",
  });

  return resp;
};
