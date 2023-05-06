import type { MeUser } from "./entities";

export interface SignupResponse {
  user: MeUser;
  token: string;
}

export interface MetaResponse {
  git_repository: string;
  git_commit: string;
  users: MetaUsers;
  members: number;
  require_invite: boolean;
}

export interface MetaUsers {
  total: number;
  active_month: number;
  active_week: number;
  active_day: number;
}

export interface UrlsResponse {
  discord?: string;
  tumblr?: string;
  google?: string;
}

export interface ExportResponse {
  path: string;
  created_at: string;
}
