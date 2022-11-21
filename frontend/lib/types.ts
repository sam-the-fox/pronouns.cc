export interface PartialPerson {
  id: string;
  name: string;
  display_name: string | null;
  avatar_urls: string[] | null;
}
export type PartialUser = PartialPerson;
export type PartialMember = PartialPerson;

interface _Person extends PartialPerson {
  bio: string | null;
  links: string[] | null;
  names: Name[] | null;
  pronouns: Pronoun[] | null;
  fields: Field[] | null;
}

export interface Member extends _Person {
  user: PartialUser;
}

export interface User extends _Person {
  members: PartialMember[] | null;
}

export type Person = Member | User;

export interface MeUser extends User {
  discord: string | null;
  discord_username: string | null;
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
  NotFound = 404,
  MethodNotAllowed = 405,
  TooManyRequests = 429,
  InternalServerError = 500,

  InvalidState = 1001,
  InvalidOAuthCode = 1002,
  InvalidToken = 1003,
  InviteRequired = 1004,
  InvalidTicket = 1005,
  InvalidUsername = 1006,
  UsernameTaken = 1007,
  InvitesDisabled = 1008,
  InviteLimitReached = 1009,
  InviteAlreadyUsed = 1010,

  UserNotFound = 2001,

  MemberNotFound = 3001,
  MemberLimitReached = 3002,

  RequestTooBig = 4001,
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
