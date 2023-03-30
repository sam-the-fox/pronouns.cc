/* eslint-disable @typescript-eslint/no-explicit-any */
import type { APIError } from "./entities";
import { PUBLIC_BASE_URL } from "$env/static/public";

export async function apiFetch<T>(
  path: string,
  {
    method,
    body,
    token,
    headers,
  }: { method?: string; body?: any; token?: string; headers?: Record<string, string> },
) {
  const resp = await fetch(`${PUBLIC_BASE_URL}/api/v1${path}`, {
    method: method || "GET",
    headers: {
      ...(token ? { Authorization: token } : {}),
      ...(headers ? headers : {}),
      "Content-Type": "application/json",
    },
    body: body ? JSON.stringify(body) : null,
  });

  const data = await resp.json();
  if (resp.status < 200 || resp.status >= 400) throw data as APIError;
  return data as T;
}

export const apiFetchClient = async <T>(path: string, method = "GET", body: any = null) =>
  apiFetch<T>(path, { method, body, token: localStorage.getItem("pronouns-token") || undefined });

/** Fetches the specified path without parsing the response body. */
export async function fastFetch(
  path: string,
  {
    method,
    body,
    token,
    headers,
  }: { method?: string; body?: any; token?: string; headers?: Record<string, string> },
) {
  const resp = await fetch(`${PUBLIC_BASE_URL}/api/v1${path}`, {
    method: method || "GET",
    headers: {
      ...(token ? { Authorization: token } : {}),
      ...(headers ? headers : {}),
      "Content-Type": "application/json",
    },
    body: body ? JSON.stringify(body) : null,
  });

  if (resp.status < 200 || resp.status >= 400) throw (await resp.json()) as APIError;
}

/** Fetches the specified path without parsing the response body. */
export const fastFetchClient = async (path: string, method = "GET", body: any = null) =>
  fastFetch(path, { method, body, token: localStorage.getItem("pronouns-token") || undefined });
