import { ErrorCode, type APIError, type Report } from "$lib/api/entities";
import { apiFetchClient } from "$lib/api/fetch";
import { error } from "@sveltejs/kit";

export const load = async () => {
  try {
    const reports = await apiFetchClient<Report[]>("/admin/reports");
    return { page: 0, isClosed: false, userId: null, reporterId: null, reports } as PageLoadData;
  } catch (e) {
    if ((e as APIError).code === ErrorCode.Forbidden) {
      throw error(400, "You're not an admin");
    }
    throw e;
  }
};

interface PageLoadData {
  page: number;
  isClosed: boolean;
  userId: string | null;
  reporterId: string | null;
  reports: Report[];
}
