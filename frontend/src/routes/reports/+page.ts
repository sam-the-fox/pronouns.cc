import { ErrorCode, type APIError, type Report } from "$lib/api/entities";
import { apiFetchClient } from "$lib/api/fetch";
import { error } from "@sveltejs/kit";

export const load = async ({ url }) => {
  const { searchParams } = url;

  const before = +(searchParams.get("before") || 0);
  const userId = searchParams.get("user_id");
  const reporterId = searchParams.get("reporter_id");
  const isClosed = searchParams.get("closed") === "true";

  try {
    let reports: Report[];
    if (userId) {
      const params = new URLSearchParams();
      if (before) params.append("before", before.toString());

      reports = await apiFetchClient<Report[]>(
        `/admin/reports/by-user/${userId}?${params.toString()}`,
      );
    } else if (reporterId) {
      const params = new URLSearchParams();
      if (before) params.append("before", before.toString());

      reports = await apiFetchClient<Report[]>(
        `/admin/reports/by-reporter/${reporterId}?${params.toString()}`,
      );
    } else {
      const params = new URLSearchParams();
      if (before) params.append("before", before.toString());
      if (isClosed) params.append("closed", "true");

      reports = await apiFetchClient<Report[]>(`/admin/reports?${params.toString()}`);
    }

    return { before, isClosed, userId, reporterId, reports } as PageLoadData;
  } catch (e) {
    if ((e as APIError).code === ErrorCode.Forbidden) {
      throw error(400, "You're not an admin");
    }
    throw e;
  }
};

interface PageLoadData {
  before: number;
  isClosed: boolean;
  userId: string | null;
  reporterId: string | null;
  reports: Report[];
}
