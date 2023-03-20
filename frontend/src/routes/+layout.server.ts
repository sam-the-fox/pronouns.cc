import { error } from "@sveltejs/kit";
import type { LayoutServerLoad } from "./$types";
import type { APIError } from "$lib/api/entities";
import { apiFetch } from "$lib/api/fetch";
import type { MetaResponse } from "$lib/api/responses";

export const load = (async (event) => {
  try {
    return await apiFetch<MetaResponse>("/meta", {});
  } catch (e) {
    throw error(500, (e as APIError).message);
  }
}) satisfies LayoutServerLoad;
