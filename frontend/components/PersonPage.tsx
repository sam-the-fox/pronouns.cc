import Head from "next/head";
import {
  Field,
  Member,
  Name,
  PartialPerson,
  Person,
  Pronoun,
  User,
  WordStatus,
} from "../lib/types";
import ReactMarkdown from "react-markdown";
import { userState } from "../lib/state";
import { useRecoilValue } from "recoil";
import FallbackImage from "./FallbackImage";
import {
  EmojiLaughing,
  HandThumbsDown,
  HandThumbsUp,
  HeartFill,
  People,
} from "react-bootstrap-icons";
import BlueLink from "./BlueLink";
import React from "react";
import Card from "./Card";

export default function PersonPage({ person }: { person: Person }) {
  return (
    <>
      <Head>
        <title key="title">{`${personFullHandle(person)} - pronouns.cc`}</title>
      </Head>
      <PersonHead person={person} />
      <IsOwnUserPageNotice person={person} />
      <div className="container mx-auto pb-[20vh]">
        <div
          className="
          m-2 p-2
          flex flex-col lg:flex-row
          justify-center lg:justify-start
          items-center lg:items-start
          lg:space-x-16
          space-y-4 lg:space-y-0
          border-b border-slate-200 dark:border-slate-700
        "
        >
          <PersonAvatar person={person} />
          <PersonInfo person={person} />
        </div>
        <LabelList content={person.names ?? []} />
        <LabelList content={person.pronouns ?? []} />
        <FieldCardGrid fields={person.fields ?? []} />
        {"user" in person ? (
          <BlueLink to={personURL(person.user)}>{`< ${
            person.user.display_name ?? person.user.name
          }`}</BlueLink>
        ) : (
          <MemberList user={person as any as User} />
        )}
      </div>
    </>
  );
}

/** Full handle of a person. */
function personFullHandle(person: Person) {
  return "user" in person
    ? `&${person.name}@${person.user.name}`
    : `@${person.name}`;
}

/** Short handle of a person. */
function personShortHandle(person: Person) {
  return ("user" in person ? "&" : "@") + person.name;
}

/** The user (account) associated with a person. */
function personUser(person: Person) {
  return "user" in person ? person.user : person;
}

/** The (relative) URL pointing to a person. */
function personURL(person: PartialPerson | Member) {
  const domain =
    typeof window !== "undefined" ? window.location.origin : process.env.DOMAIN;
  return `${domain}/u/${"user" in person ? `${person.user.name}/` : ""}${
    person.name
  }`;
}

function PersonHead({ person }: { person: Person }) {
  const {
    id,
    name,
    display_name,
    avatar_urls,
    bio,
    links,
    names,
    pronouns,
    fields,
  } = person;
  let description = "";
  if (
    names?.filter((name) => name.status === WordStatus.Favourite)?.length &&
    pronouns?.filter((pronoun) => pronoun.status === WordStatus.Favourite)
      ?.length
  ) {
    description = `${personShortHandle(person)} goes by ${names
      .filter((name) => name.status === WordStatus.Favourite)
      .map((name) => name.name)
      .join(", ")} and uses ${pronouns
      .filter((pronoun) => pronoun.status === WordStatus.Favourite)
      .map(
        (pronoun) =>
          pronoun.display_text ??
          pronoun.pronouns.split("/").slice(0, 2).join("/")
      )
      .join(", ")} pronouns.`;
  } else if (
    names?.filter((name) => name.status === WordStatus.Favourite)?.length
  ) {
    description = `${personShortHandle(person)} goes by ${names
      .filter((name) => name.status === WordStatus.Favourite)
      .map((name) => name.name)
      .join(", ")}.`;
  } else if (
    pronouns?.filter((pronoun) => pronoun.status === WordStatus.Favourite)
      ?.length
  ) {
    description = `${personShortHandle(person)} uses ${pronouns
      .filter((pronoun) => pronoun.status === WordStatus.Favourite)
      .map(
        (pronoun) =>
          pronoun.display_text ??
          pronoun.pronouns.split("/").slice(0, 2).join("/")
      )
      .join(", ")} pronouns.`;
  } else if (bio && bio !== "") {
    description = bio.slice(0, 500);
  }

  return (
    <Head>
      <meta key="og:sitename" property="og:site_name" content="pronouns.cc" />
      <meta
        key="og:title"
        property="og:title"
        content={
          display_name
            ? `${display_name} (${personFullHandle(person)})`
            : personFullHandle(person)
        }
      />
      {avatar_urls && avatar_urls.length > 0 ? (
        <meta key="og:image" property="og:image" content={avatar_urls[0]} />
      ) : (
        <></>
      )}
      <meta
        key="og:description"
        property="og:description"
        content={description}
      />
      <meta key="og:url" property="og:url" content={personURL(person)} />
    </Head>
  );
}

function IsOwnUserPageNotice({ person }: { person: Person }) {
  return useRecoilValue(userState)?.id === person.id ? (
    <div className="lg:w-1/3 mx-auto bg-slate-100 dark:bg-slate-700 shadow rounded-md p-2">
      You are currently viewing your <b>public</b> user profile.
      <br />
      <BlueLink to="/edit/profile">Edit your profile</BlueLink>
    </div>
  ) : (
    <></>
  );
}

function MemberList({ user, className }: { user: User; className?: string }) {
  const partialMembers = user.members;
  return (
    <div className={`mx-auto flex-col items-center ${className || ""}`}>
      <h1 className="text-2xl">Members</h1>
      <ul>
        {partialMembers?.map((partialMember) => (
          <li className='before:[content:"-_"]' key={partialMember.id}>
            <BlueLink to={`/u/${user.name}/${partialMember.name}`}>
              <span>{partialMember.display_name ?? partialMember.name}</span>
            </BlueLink>
          </li>
        ))}
      </ul>
    </div>
  );
}

function PersonAvatar({ person }: { person: Person }) {
  const { display_name, name, avatar_urls } = person;
  return avatar_urls && avatar_urls.length !== 0 ? (
    <FallbackImage
      className="max-w-xs rounded-full"
      urls={avatar_urls}
      alt={`${display_name || name}'s avatar`}
    />
  ) : (
    <></>
  );
}

function PersonInfo({ person }: { person: Person }) {
  const { display_name, name, bio, links } = person;
  return (
    <div className="flex flex-col">
      {/* name */}
      <h1 className="text-2xl font-bold">
        {display_name === null ? name : display_name}
      </h1>
      {/* handle */}
      <h3 className="text-xl font-light text-slate-600 dark:text-slate-400">
        {personFullHandle(person)}
      </h3>
      {/* bio */}
      {bio && (
        <ReactMarkdown className="prose dark:prose-invert prose-slate">
          {bio}
        </ReactMarkdown>
      )}
      {/* links */}
      {links?.length && (
        <div className="flex flex-col mx-auto lg:ml-auto">
          {links.map((link, index) => (
            <a
              key={index}
              href={link}
              rel="nofollow noopener noreferrer me"
              className="hover:underline text-sky-500 dark:text-sky-400"
            >
              {link}
            </a>
          ))}
        </div>
      )}
    </div>
  );
}

function LabelList({ content }: { content: Name[] | Pronoun[] }) {
  return content?.length > 0 ? (
    <div className="border-b border-slate-200 dark:border-slate-700">
      {content.map((label, index) => (
        <LabelLine key={index} label={label} />
      ))}
    </div>
  ) : (
    <></>
  );
}

function LabelStatusIcon({ status }: { status: WordStatus }) {
  return React.createElement(
    {
      [WordStatus.Favourite]: HeartFill,
      [WordStatus.Okay]: HandThumbsUp,
      [WordStatus.Jokingly]: EmojiLaughing,
      [WordStatus.FriendsOnly]: People,
      [WordStatus.Avoid]: HandThumbsDown,
    }[status],
    { className: "inline" }
  );
}

function LabelsLine({ labels }: { labels: Name[] | Pronoun[] }) {
  if (!labels?.length) return <></>;
  const status = labels[0].status;
  const text = labels
    .map((label) =>
      "name" in label
        ? label.name
        : label.display_text ?? label.pronouns.split("/").slice(0, 2).join("/")
    )
    .join(", ");
  return (
    <p
      className={`
      ${status === WordStatus.Favourite ? "text-lg font-bold" : ""}
      ${
        status === WordStatus.Avoid ? "text-slate-600 dark:text-slate-400" : ""
      }`}
    >
      <LabelStatusIcon status={status} /> {text}
    </p>
  );
}

function LabelLine({ label }: { label: Name | Pronoun }) {
  return <LabelsLine labels={[label] as Name[] | Pronoun[]} />;
}

function FieldCardGrid({ fields }: { fields: Field[] }) {
  return (
    <div className="flex flex-col md:flex-row gap-4 py-2 [&>*]:flex-1">
      {fields?.map((field, index) => (
        <FieldCard field={field} key={index} />
      ))}
    </div>
  );
}

const fieldEntryStatus: { [key in string]: WordStatus } = {
  favourite: WordStatus.Favourite,
  okay: WordStatus.Okay,
  jokingly: WordStatus.Jokingly,
  friends_only: WordStatus.FriendsOnly,
  avoid: WordStatus.Avoid,
};

function FieldCard({
  field,
  draggable,
}: {
  field: Field;
  draggable?: boolean;
}) {
  return (
    <Card title={field.name} draggable={draggable}>
      {Object.entries(fieldEntryStatus).map(([statusName, status], i) => (
        <LabelsLine
          key={i}
          labels={(field as any)[statusName]?.map((name: string) => ({
            name,
            status,
          }))}
        />
      ))}
    </Card>
  );
}
