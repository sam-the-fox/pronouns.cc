import * as API from "./api-fetch";
import { fetchAPI } from "./api-fetch";

function getDomain(): string {
  const domain =
    typeof window !== "undefined" ? window.location.origin : process.env.DOMAIN;
  if (!domain) throw new Error("process.env.DOMAIN not set");
  return domain;
}

export class PartialPerson {
  id: string;
  name: string;
  displayName: string | null;
  avatarUrls: string[];
  constructor({ id, name, display_name, avatar_urls }: API.PartialPerson) {
    this.id = id;
    this.name = name;
    this.displayName = display_name;
    this.avatarUrls = avatar_urls ?? [];
  }

  display(): string {
    return this.displayName ?? this.name;
  }
}

export class PartialUser extends PartialPerson {}

export class PartialMember extends PartialPerson {}

abstract class _Person extends PartialPerson {
  bio: string | null;
  links: string[];
  names: Name[];
  pronouns: Pronouns[];
  fields: Field[];
  constructor(apiData: API._Person) {
    super(apiData);
    const { bio, links, names, pronouns, fields } = apiData;
    this.bio = bio;
    this.links = links ?? [];
    this.names = (names ?? []).map((x) => new Name(x));
    this.pronouns = (pronouns ?? []).map((x) => new Pronouns(x));
    this.fields = (fields ?? []).map((x) => new Field(x));
  }

  abstract fullHandle(): string;

  shortHandle(): string {
    return this.fullHandle();
  }

  abstract relativeURL(): string;

  absoluteURL(): string {
    return `${getDomain()}${this.relativeURL()}`;
  }
}

export class User extends _Person {
  partialMembers: PartialMember[];
  constructor(apiData: API.User) {
    super(apiData);
    const { members } = apiData;
    this.partialMembers = (members ?? []).map((x) => new PartialMember(x));
  }

  static async fetchFromName(name: string): Promise<User> {
    return new User(await fetchAPI<API.User>(`/users/${name}`));
  }

  fullHandle(): string {
    return `@${this.name}`;
  }

  shortHandle(): string {
    return this.fullHandle();
  }

  relativeURL(): string {
    return `/u/${this.name}`;
  }
}

export class Member extends _Person {
  partialUser: PartialUser;
  constructor(apiData: API.Member) {
    super(apiData);
    const { user } = apiData;
    this.partialUser = new PartialUser(user);
  }

  static async fetchFromUserAndMemberName(
    userName: string,
    memberName: string
  ): Promise<Member> {
    return new Member(
      await fetchAPI<API.Member>(`/users/${userName}/members/${memberName}`)
    );
  }

  fullHandle(): string {
    return `${this.name}@${this.partialUser.name}`;
  }

  relativeURL(): string {
    return `/u/${this.partialUser.name}/${this.name}`;
  }
}

export type Person = Member | User;

export class MeUser extends User {
  discord: string | null;
  discordUsername: string | null;
  constructor(apiData: API.MeUser) {
    super(apiData);
    const { discord, discord_username } = apiData;
    this.discord = discord;
    this.discordUsername = discord_username;
  }

  static async fetchMe(): Promise<MeUser> {
    return new MeUser(await fetchAPI<API.MeUser>("/users/@me"));
  }
}

export enum LabelType {
  Name = 1,
  Pronouns = 2,
  Unspecified = 3,
}

export const LabelStatus = API.WordStatus;
export type LabelStatus = API.WordStatus;

export interface LabelData {
  type?: LabelType;
  displayText: string | null;
  text: string;
  status: LabelStatus;
}

export class Label {
  type: LabelType;
  displayText: string | null;
  text: string;
  status: LabelStatus;
  constructor({ type, displayText, text, status }: LabelData) {
    this.type = type ?? LabelType.Unspecified;
    this.displayText = displayText;
    this.text = text;
    this.status = status;
  }

  display(): string {
    return this.displayText ?? this.text;
  }

  shortDisplay(): string {
    return this.display();
  }
}

export class Name extends Label {
  constructor({ name, status }: API.Name) {
    super({
      type: LabelType.Name,
      displayText: null,
      text: name,
      status,
    });
  }
}

export class Pronouns extends Label {
  constructor({ display_text, pronouns, status }: API.Pronoun) {
    super({
      type: LabelType.Pronouns,
      displayText: display_text ?? null,
      text: pronouns,
      status,
    });
  }

  get pronouns(): string[] {
    return this.text.split("/");
  }
  set pronouns(to: string[]) {
    this.text = to.join("/");
  }

  shortDisplay(): string {
    return this.displayText ?? this.pronouns.splice(0, 2).join("/");
  }
}

export class Field {
  name: string;
  labels: Label[];
  constructor({
    name,
    favourite,
    okay,
    jokingly,
    friends_only,
    avoid,
  }: API.Field) {
    this.name = name;
    function transpose(arr: API.Arr<string>, status: LabelStatus): Label[] {
      return (arr ?? []).map(
        (text) =>
          new Label({
            displayText: null,
            text,
            status,
          })
      );
    }
    this.labels = [
      ...transpose(favourite, LabelStatus.Favourite),
      ...transpose(okay, LabelStatus.Okay),
      ...transpose(jokingly, LabelStatus.Jokingly),
      ...transpose(friends_only, LabelStatus.FriendsOnly),
      ...transpose(avoid, LabelStatus.Avoid),
    ];
  }
}
