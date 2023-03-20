import { PUBLIC_BASE_URL } from "$env/static/public";
import { apiFetch } from "$lib/api/fetch";
import type { UrlsResponse } from "$lib/api/responses";

export const load = async () => {
  const resp = await apiFetch<UrlsResponse>("/auth/urls", {
    method: "POST",
    body: {
      callback_domain: PUBLIC_BASE_URL,
    },
  });

  return { urls: resp };
};
