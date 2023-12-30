import { ErrorCode, type APIError } from "$lib/api/entities";
import { apiFetchClient } from "$lib/api/fetch";
import type { ExportResponse } from "$lib/api/responses";
import { error } from "@sveltejs/kit";

export const load = async () => {
  try {
    const data = await apiFetchClient<ExportResponse>("/users/@me/export");
    return { exportData: data };
  } catch (e) {
    if ((e as APIError).code === ErrorCode.NotFound) return { exportData: null };

    error(500, e as App.Error);
  }
};
