import { apiFetch } from "$lib/api/fetch";
import { PUBLIC_BASE_URL } from "$env/static/public";

export const load = async () => {
  const resp = await apiFetch<UrlsResponse>("/auth/urls", {
    method: "POST",
    body: {
      callback_domain: PUBLIC_BASE_URL,
    },
  });

  return resp;
};

interface UrlsResponse {
  discord: string;
}
