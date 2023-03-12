export interface User {
  id: string;
  name: string;
  display_name: string | null;
  bio: string | null;
  avatar_urls: string[];
  links: string[];

  names: FieldEntry[];
  pronouns: Pronoun[];
  members: PartialMember[];
  fields: Field[];
}

export interface MeUser extends User {
  discord: string | null;
  discord_username: string | null;
}

export interface Field {
  name: string;
  entries: FieldEntry[];
}

export interface FieldEntry {
  value: string;
  status: WordStatus;
}

export interface Pronoun {
  pronouns: string;
  display_text: string | null;
  status: WordStatus;
}

export enum WordStatus {
  Unknown = 0,
  Favourite = 1,
  Okay = 2,
  Jokingly = 3,
  FriendsOnly = 4,
  Avoid = 5,
}

export interface PartialMember {
  id: string;
  name: string;
  display_name: string | null;
  bio: string | null;
  avatar_urls: string[];
  links: string[];
  names: FieldEntry[];
  pronouns: Pronoun[];
}

export interface Member extends PartialMember {
  fields: Field[];

  user: MemberPartialUser;
}

export interface MemberPartialUser {
  id: string;
  name: string;
  display_name: string | null;
  avatar_urls: string[];
}

export interface APIError {
  code: ErrorCode;
  message?: string;
  details?: string;
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
