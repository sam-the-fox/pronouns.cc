import type { MeUser } from "./entities";

export interface SignupResponse {
  user: MeUser;
  token: string;
}

export interface MetaResponse {
  git_repository: string;
  git_commit: string;
  users: number;
  members: number;
  require_invite: boolean;
}

export interface UrlsResponse {
  discord: string;
  tumblr: string;
  google: string;
}

export interface ExportResponse {
  path: string;
  created_at: string;
}
