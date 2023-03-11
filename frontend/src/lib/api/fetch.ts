import type { APIError } from "./entities";
import { PUBLIC_BASE_URL } from "$env/static/public";

export async function apiFetch<T>(
  path: string,
  { method, body, token }: { method?: string; body?: any; token?: string },
) {

  const resp = await fetch(`${PUBLIC_BASE_URL}/api/v1${path}`, {
    method: method || "GET",
    headers: {
      ...(token ? { Authorization: token } : {}),
      "Content-Type": "application/json",
    },
    body: body ? JSON.stringify(body) : null,
  });

  const data = await resp.json();
  if (resp.status < 200 || resp.status >= 300) throw data as APIError;
  return data as T;
}

export const apiFetchClient = async <T>(path: string, method: string = "GET", body: any = null) =>
  apiFetch<T>(path, { method, body, token: localStorage.getItem("pronouns-token") || undefined });
