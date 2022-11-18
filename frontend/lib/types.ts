export interface MeUser extends User {
  discord: string | null;
  discord_username: string | null;
}

export interface User {
  id: string;
  username: string;
  display_name: string | null;
  bio: string | null;
  avatar_urls: string[] | null;
  links: string[] | null;
  members: PartialMember[];
  names: Name[];
  pronouns: Pronoun[];
  fields: Field[];
}

export interface PartialMember {
  id: string;
  name: string;
  avatar_urls: string[] | null;
}

export interface Member extends PartialMember {
  bio: string | null;
  links: string[] | null;
  id: string;
  name: string;
  avatar_urls: string[] | null;

  user?: PartialUser;
}

export interface Name {
  name: string;
  status: WordStatus;
}

export interface Pronoun {
  display_text?: string;
  pronouns: string;
  status: WordStatus;
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

export enum WordStatus {
  Favourite = 1,
  Okay = 2,
  Jokingly = 3,
  FriendsOnly = 4,
  Avoid = 5,
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

export interface SignupResponse {
  user: MeUser;
  token: string;
}

export interface PartialUser {
  id: string;
  username: string;
  display_name: string | null;
  avatar_urls: string[] | null;
}
