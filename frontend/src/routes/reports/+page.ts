import type { Report } from "$lib/api/entities";
import { apiFetchClient } from "$lib/api/fetch";

export const load = async () => {
  const reports = await apiFetchClient<Report[]>("/admin/reports");
  return { page: 0, isClosed: false, userId: null, reporterId: null, reports } as PageLoadData;
};

interface PageLoadData {
  page: number;
  isClosed: boolean;
  userId: string | null;
  reporterId: string | null;
  reports: Report[];
}
