import { apiFetchClient } from "$lib/api/fetch";

export const load = async () => {
  const tokens = await apiFetchClient<Token[]>("/auth/tokens");
  return { tokens };
};

interface Token {
  id: string;
  created: string;
  expires: string;
}
