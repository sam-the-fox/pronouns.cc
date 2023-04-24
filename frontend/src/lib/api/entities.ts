import { PUBLIC_BASE_URL } from "$env/static/public";

export const MAX_MEMBERS = 500;
export const MAX_DESCRIPTION_LENGTH = 1000;

export interface User {
  id: string;
  name: string;
  display_name: string | null;
  bio: string | null;
  avatar: string | null;
  links: string[];
  member_title: string | null;

  names: FieldEntry[];
  pronouns: Pronoun[];
  members: PartialMember[];
  fields: Field[];
  custom_preferences: CustomPreferences;
}

export interface CustomPreferences {
  [key: string]: CustomPreference;
}

export interface CustomPreference {
  icon: string;
  tooltip: string;
  size: PreferenceSize;
  muted: boolean;
  favourite: boolean;
}

export enum PreferenceSize {
  Large = "large",
  Normal = "normal",
  Small = "small",
}

export interface MeUser extends User {
  max_invites: number;
  discord: string | null;
  discord_username: string | null;
  tumblr: string | null;
  tumblr_username: string | null;
  google: string | null;
  google_username: string | null;
  fediverse: string | null;
  fediverse_username: string | null;
  fediverse_instance: string | null;
  list_private: boolean;
}

export interface Field {
  name: string;
  entries: FieldEntry[];
}

export interface FieldEntry {
  value: string;
  status: string;
}

export interface Pronoun {
  pronouns: string;
  display_text: string | null;
  status: string;
}

export interface PartialMember {
  id: string;
  name: string;
  display_name: string | null;
  bio: string | null;
  avatar: string | null;
  links: string[];
  names: FieldEntry[];
  pronouns: Pronoun[];
}

export interface Member extends PartialMember {
  fields: Field[];

  user: MemberPartialUser;
  unlisted?: boolean;
}

export interface MemberPartialUser {
  id: string;
  name: string;
  display_name: string | null;
  avatar: string | null;
  custom_preferences: CustomPreferences;
}

export interface Invite {
  code: string;
  created: string;
  used: boolean;
}

export interface Report {
  id: string;
  user_id: string;
  user_name: string;
  member_id: string | null;
  member_name: string | null;
  reason: string;
  reporter_id: string;

  created_at: string;
  resolved_at: string | null;
  admin_id: string | null;
  admin_comment: string | null;
}

export interface Warning {
  id: number;
  reason: string;
  created_at: string;
  read: boolean;
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
  RecentExport = 1012,
  UnsupportedInstance = 1013,
  AlreadyLinked = 1014,
  NotLinked = 1015,
  LastProvider = 1016,
  InvalidCaptcha = 1017,

  UserNotFound = 2001,

  MemberNotFound = 3001,
  MemberLimitReached = 3002,
  MemberNameInUse = 3003,
  NotOwnMember = 3004,

  RequestTooBig = 4001,
  MissingPermissions = 4002,
}

export const pronounDisplay = (entry: Pronoun) => {
  if (entry.display_text) return entry.display_text;

  const split = entry.pronouns.split("/");
  if (split.length < 2) return split.join("/");
  else return split.slice(0, 2).join("/");
};

export const userAvatars = (user: User | MeUser | MemberPartialUser) => {
  if (!user.avatar) return defaultAvatars;

  return [
    `${PUBLIC_BASE_URL}/media/users/${user.id}/${user.avatar}.webp`,
    `${PUBLIC_BASE_URL}/media/users/${user.id}/${user.avatar}.jpg`,
  ];
};

export const memberAvatars = (member: Member | PartialMember) => {
  if (!member.avatar) return defaultAvatars;

  return [
    `${PUBLIC_BASE_URL}/media/members/${member.id}/${member.avatar}.webp`,
    `${PUBLIC_BASE_URL}/media/members/${member.id}/${member.avatar}.jpg`,
  ];
};

export const defaultAvatars = [
  `${PUBLIC_BASE_URL}/default/512.webp`,
  `${PUBLIC_BASE_URL}/default/512.jpg`,
];

export interface PronounsJson {
  pages: Pronouns;
  autocomplete: Pronouns;
}

interface Pronouns {
  [key: string]: { pronouns: string[]; display?: string };
}
