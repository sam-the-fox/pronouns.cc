/* eslint-disable @typescript-eslint/no-explicit-any */
import { ErrorCode, type APIError } from "./entities";
import { PUBLIC_BASE_URL } from "$env/static/public";
import { addToast } from "$lib/toast";
import { userStore } from "$lib/store";

export async function apiFetch<T>(
  path: string,
  {
    method,
    body,
    token,
    headers,
    version,
  }: {
    method?: string;
    body?: any;
    token?: string;
    headers?: Record<string, string>;
    version?: number;
  },
) {
  const resp = await fetch(`${PUBLIC_BASE_URL}/api/v${version || 1}${path}`, {
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

export const apiFetchClient = async <T>(
  path: string,
  method = "GET",
  body: any = null,
  version = 1,
) => {
  try {
    const data = await apiFetch<T>(path, {
      method,
      body,
      token: localStorage.getItem("pronouns-token") || undefined,
      version,
    });
    return data;
  } catch (e) {
    if ((e as APIError).code === ErrorCode.InvalidToken) {
      addToast({ header: "Token expired", body: "Your token has expired, please log in again." });
      userStore.set(null);
      localStorage.removeItem("pronouns-token");
      localStorage.removeItem("pronouns-user");
    }
    throw e;
  }
};

/** Fetches the specified path without parsing the response body. */
export async function fastFetch(
  path: string,
  {
    method,
    body,
    token,
    headers,
    version,
  }: {
    method?: string;
    body?: any;
    token?: string;
    headers?: Record<string, string>;
    version?: number;
  },
) {
  const resp = await fetch(`${PUBLIC_BASE_URL}/api/v${version || 1}${path}`, {
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
export const fastFetchClient = async (
  path: string,
  method = "GET",
  body: any = null,
  version = 1,
) => {
  try {
    await fastFetch(path, {
      method,
      body,
      token: localStorage.getItem("pronouns-token") || undefined,
      version,
    });
  } catch (e) {
    if ((e as APIError).code === ErrorCode.InvalidToken) {
      addToast({ header: "Token expired", body: "Your token has expired, please log in again." });
      userStore.set(null);
      localStorage.removeItem("pronouns-token");
      localStorage.removeItem("pronouns-user");
    }
    throw e;
  }
};
