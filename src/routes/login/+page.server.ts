import { fetchAPI } from "$lib/api/fetch";

export const load = async () => {
  const resp = await fetchAPI<UrlsResponse>("/auth/urls", {
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
