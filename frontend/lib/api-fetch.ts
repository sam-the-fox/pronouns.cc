/** An array returned by the API. (Can be `null` due to a quirk in Go.) */
export type Arr<T> = T[] | null;

export interface PartialPerson {
  id: string;
  name: string;
  display_name: string | null;
  avatar_urls: Arr<string>;
}

export type PartialUser = PartialPerson;

export type PartialMember = PartialPerson;

/** The shared interface of `Member` and `User`.
 * A typical `_Person` is only one of those two, so consider using `Person` instead.
 */
export interface _Person extends PartialPerson {
  bio: string | null;
  links: Arr<string>;
  names: Arr<FieldEntry>;
  pronouns: Arr<Pronoun>;
  fields: Arr<Field>;
}

export interface User extends _Person {
  members: Arr<PartialMember>;
}

export interface Member extends _Person {
  user: PartialUser;
}

export type Person = Member | User;

export interface MeUser extends User {
  discord: string | null;
  discord_username: string | null;
}

export interface Pronoun {
  display_text?: string;
  pronouns: string;
  status: WordStatus;
}

export interface Field {
  name: string;
  entries: Arr<FieldEntry>;
}

export interface FieldEntry {
  value: string;
  status: WordStatus;
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

const apiBase = process.env.API_BASE ?? "/api";

export async function fetchAPI<T>(
  path: string,
  method = "GET",
  body: any = null
) {
  const token =
    typeof localStorage !== "undefined" &&
    localStorage.getItem("pronouns-token");

  const resp = await fetch(`${apiBase}/v1${path}`, {
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
