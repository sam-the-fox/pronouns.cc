export interface User {
  id: string;
  username: string;
  display_name: string | null;
  bio: string | null;
  avatar_urls: string[] | null;
  links: string[] | null;
}

export interface MeUser extends User {
  discord: string | null;
  discord_username: string | null;
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
