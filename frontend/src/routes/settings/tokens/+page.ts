import { apiFetchClient } from "$lib/api/fetch";

export const load = async () => {
  const tokens = await apiFetchClient<Token[]>("/auth/tokens");
  return { tokens: tokens.filter((token) => token.api_only) };
};

interface Token {
  id: string;
  api_only: boolean;
  read_only: boolean;
  created: string;
  expires: string;
}
