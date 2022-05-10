export interface MeUser {
  id: string;
  username: string;
  display_name: string | null;
  bio: string | null;
  avatar_source: string | null;
  avatar_url: string | null;
  links: string[] | null;
  discord: string | null;
  discord_username: string | null;
}

export interface User {
  id: string;
  username: string;
  display_name: string | null;
  bio: string | null;
  avatar_source: string | null;
  links: string[] | null;
  members: PartialMember[];
}

export interface PartialMember {
  id: string;
  name: string;
  avatar_url: string | null;
}

export interface APIError {
  code: ErrorCode;
  message?: string;
}

export enum ErrorCode {
  BadRequest = 400,
  Forbidden = 403,
  InternalServerError = 500,

  InvalidState = 1001,
  InvalidOAuthCode = 1002,

  UserNotFound = 2001,
}
