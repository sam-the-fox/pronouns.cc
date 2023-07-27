import { building } from "$app/environment";
import type { LayoutServerLoad } from "./$types";
import { apiFetch } from "$lib/api/fetch";
import type { MetaResponse } from "$lib/api/responses";

export const load = (async () => {
  try {
    return await apiFetch<MetaResponse>("/meta", {});
  } catch (e) {
    console.warn("error fetching meta endpoint:", e);

    if (building) {
      // just return an empty object--this only affects the three static pages, nothing else, so it's fine
      return {
        git_repository: "",
        git_commit: "",
        users: {
          total: 0,
          active_month: 0,
          active_week: 0,
          active_day: 0,
        },
        members: 0,
        require_invite: false,
      };
    } else {
      throw e;
    }
  }
}) satisfies LayoutServerLoad;
