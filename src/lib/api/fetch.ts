import type { APIError } from "./entities";

export async function fetchAPI<T>(
  path: string,
  { method, body, token }: { method?: string; body?: any; token?: string },
) {
  const resp = await fetch(`${process.env.ORIGIN}/api/v1${path}`, {
    method,
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
