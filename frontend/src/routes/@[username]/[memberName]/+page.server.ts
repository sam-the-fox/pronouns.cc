import { apiFetch } from "$lib/api/fetch";
import type { Member } from "$lib/api/entities";

export const load = async ({ params }) => {
  const resp = await apiFetch<Member>(`/users/${params.username}/members/${params.memberName}`, {
    method: "GET",
  });

  return resp;
};
