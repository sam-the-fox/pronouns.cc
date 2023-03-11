import { apiFetch } from "$lib/api/fetch";

export const load = async () => {
  const resp = await apiFetch<UrlsResponse>("/auth/urls", {
    method: "POST",
    body: {
      callback_domain: process.env.ORIGIN,
    },
  });

  return resp;
};

interface UrlsResponse {
  discord: string;
}
