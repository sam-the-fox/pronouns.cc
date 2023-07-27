import { apiFetchClient } from "$lib/api/fetch";
import type { PrideFlag } from "$lib/api/entities";

export const load = async () => {
  const data = await apiFetchClient<PrideFlag[]>("/users/@me/flags");
  return { flags: data };
};
