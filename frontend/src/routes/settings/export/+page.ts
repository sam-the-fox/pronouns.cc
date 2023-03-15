import { ErrorCode, type APIError } from "$lib/api/entities";
import { apiFetchClient } from "$lib/api/fetch";
import { error } from "@sveltejs/kit";

export const load = async () => {
  try {
    const data = await apiFetchClient<ExportResponse>("/auth/export");
    return { exportData: data };
  } catch (e) {
    if ((e as APIError).code === ErrorCode.NotFound) return { exportData: null };

    throw error((e as APIError).code, (e as APIError).message);
  }
};

interface ExportResponse {
  path: string;
  created_at: Date;
}
