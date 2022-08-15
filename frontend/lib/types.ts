export interface MeUser extends User {
  avatar_source: string | null;
  discord: string | null;
  discord_username: string | null;
}

export interface User {
  id: string;
  username: string;
  display_name: string | null;
  bio: string | null;
  avatar_url: string | null;
  links: string[] | null;
  members: PartialMember[];
  fields: Field[];
}

export interface PartialMember {
  id: string;
  name: string;
  avatar_url: string | null;
}

export interface Field {
  name: string;
  favourite: string[] | null;
  okay: string[] | null;
  jokingly: string[] | null;
  friends_only: string[] | null;
  avoid: string[] | null;
}

export interface APIError {
  code: ErrorCode;
  message?: string;
  details?: string;
}

export enum ErrorCode {
  BadRequest = 400,
  Forbidden = 403,
  InternalServerError = 500,

  InvalidState = 1001,
  InvalidOAuthCode = 1002,
  InvalidToken = 1003,

  UserNotFound = 2001,
}

export interface SignupRequest {
  username: string;
  ticket: string;
  invite_code?: string;
}
