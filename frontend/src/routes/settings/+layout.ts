import type { LayoutLoad } from "./$types";

export const ssr = false;

export const load = (async ({ parent }) => {
  const data = await parent();
  return data;
}) satisfies LayoutLoad;
