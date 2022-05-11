import axios from "axios";
import type { APIError } from "./types";

export default async function fetchAPI<T>(path: string) {
  const resp = await axios.get<T | APIError>(`/api/v1${path}`, {
    headers: {
      Authorization: localStorage.getItem("pronouns-token"),
      "Content-Type": "application/json",
    },
  });
  if (resp.status !== 200) throw resp.data as APIError;

  return resp.data as T;
}
