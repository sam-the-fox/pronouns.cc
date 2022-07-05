import axios from "axios";
import type { APIError } from "./types";

export default async function fetchAPI<T>(
  path: string,
  method = "GET",
  body: any = null
) {
  let headers = {};
  const token = localStorage.getItem("pronouns-token");
  if (token) {
    headers = {
      Authorization: token,
    };
  }

  const resp = await fetch(`/api/v1${path}`, {
    method,
    headers: {
      ...headers,
      "Content-Type": "application/json",
    },
    body: body ? JSON.stringify(body) : null,
  });

  const data = await resp.json();
  if (resp.status !== 200) throw data as APIError;
  return data as T;
}
