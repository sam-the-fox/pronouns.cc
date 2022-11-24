import Head from "next/head";
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
import { Field, Label, LabelStatus, Person, User } from "../lib/api";

export default function PersonPage({ person }: { person: Person }) {
  return (
    <>
      <Head>
        <title key="title">{`${person.fullHandle()} - pronouns.cc`}</title>
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
        <LabelList content={person.names} />
        <LabelList content={person.pronouns} />
        <FieldCardGrid fields={person.fields} />
        {person instanceof User ? (
          <MemberList user={person} />
        ) : (
          <BlueLink
            to={person.relativeURL()}
          >{`< ${person.display()}`}</BlueLink>
        )}
      </div>
    </>
  );
}

function PersonHead({ person }: { person: Person }) {
  const { displayName, avatarUrls, bio, names, pronouns } = person;
  let description = "";
  const favNames = names.filter((x) => x.status === LabelStatus.Favourite);
  const favPronouns = pronouns.filter(
    (x) => x.status === LabelStatus.Favourite
  );
  if (favNames.length || favPronouns.length) {
    description = `${person.shortHandle()}${
      favNames.length
        ? ` goes by ${favNames.map((x) => x.display()).join(", ")}`
        : ""
    }${favNames.length && favPronouns.length ? " and" : ""}${
      favPronouns.length
        ? `uses ${favPronouns
            .map((x) => x.shortDisplay())
            .join(", ")} pronouns.`
        : ""
    }`;
  } else if (bio && bio !== "") {
    description = `${bio.slice(0, 500)}${bio.length > 500 ? "…" : ""}`;
  }

  return (
    <Head>
      <meta key="og:sitename" property="og:site_name" content="pronouns.cc" />
      <meta
        key="og:title"
        property="og:title"
        content={
          displayName
            ? `${displayName} (${person.fullHandle()})`
            : person.fullHandle()
        }
      />
      {avatarUrls && avatarUrls.length > 0 && (
        <meta key="og:image" property="og:image" content={avatarUrls[0]} />
      )}
      <meta
        key="og:description"
        property="og:description"
        content={description}
      />
      <meta key="og:url" property="og:url" content={person.absoluteURL()} />
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
  const partialMembers = user.partialMembers;
  return (
    <div className={`mx-auto flex-col items-center ${className || ""}`}>
      <h1 className="text-2xl">Members</h1>
      <ul>
        {partialMembers.map((partialMember) => (
          <li className='before:[content:"-_"]' key={partialMember.id}>
            <BlueLink to={`/u/${user.name}/${partialMember.name}`}>
              <span>{partialMember.displayName ?? partialMember.name}</span>
            </BlueLink>
          </li>
        ))}
      </ul>
    </div>
  );
}

function PersonAvatar({ person }: { person: Person }) {
  const { displayName, name, avatarUrls } = person;
  return avatarUrls && avatarUrls.length !== 0 ? (
    <FallbackImage
      className="max-w-xs rounded-full"
      urls={avatarUrls}
      alt={`${displayName || name}'s avatar`}
    />
  ) : (
    <></>
  );
}

function PersonInfo({ person }: { person: Person }) {
  const { displayName, name, bio, links } = person;
  return (
    <div className="flex flex-col">
      {/* name */}
      <h1 className="text-2xl font-bold">
        {displayName === null ? name : displayName}
      </h1>
      {/* handle */}
      <h3 className="text-xl font-light text-slate-600 dark:text-slate-400">
        {person.fullHandle()}
      </h3>
      {/* bio */}
      {bio && (
        <ReactMarkdown className="prose dark:prose-invert prose-slate">
          {bio}
        </ReactMarkdown>
      )}
      {/* links */}
      {links.length > 0 && (
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

function LabelList({ content }: { content: Label[] }) {
  return content.length > 0 ? (
    <div className="border-b border-slate-200 dark:border-slate-700">
      {content.map((label, index) => (
        <LabelLine key={index} label={label} />
      ))}
    </div>
  ) : (
    <></>
  );
}

function LabelStatusIcon({ status }: { status: LabelStatus }) {
  return React.createElement(
    {
      [LabelStatus.Favourite]: HeartFill,
      [LabelStatus.Okay]: HandThumbsUp,
      [LabelStatus.Jokingly]: EmojiLaughing,
      [LabelStatus.FriendsOnly]: People,
      [LabelStatus.Avoid]: HandThumbsDown,
    }[status],
    { className: "inline" }
  );
}

function LabelsLine({
  status,
  texts,
}: {
  status: LabelStatus;
  texts: string[];
}) {
  return !texts.length ? (
    <></>
  ) : (
    <p
      className={`
      ${status === LabelStatus.Favourite ? "text-lg font-bold" : ""}
      ${
        status === LabelStatus.Avoid ? "text-slate-600 dark:text-slate-400" : ""
      }`}
    >
      <LabelStatusIcon status={status} /> {texts.join(", ")}
    </p>
  );
}

function LabelLine({ label }: { label: Label }) {
  return <LabelsLine status={label.status} texts={[label.display()]} />;
}

function FieldCardGrid({ fields }: { fields: Field[] }) {
  return (
    <div className="flex flex-col md:flex-row gap-4 py-2 [&>*]:flex-1">
      {fields.map((field, index) => (
        <FieldCard field={field} key={index} />
      ))}
    </div>
  );
}

const labelStatusOrder: LabelStatus[] = [
  LabelStatus.Favourite,
  LabelStatus.Okay,
  LabelStatus.Jokingly,
  LabelStatus.FriendsOnly,
  LabelStatus.Avoid,
];

function FieldCard({
  field,
  draggable,
}: {
  field: Field;
  draggable?: boolean;
}) {
  return (
    <Card title={field.name} draggable={draggable}>
      {labelStatusOrder.map((status, i) => (
        <LabelsLine
          key={i}
          status={status}
          texts={field.labels
            .filter((x) => x.status === status)
            .map((x) => x.display())}
        />
      ))}
    </Card>
  );
}
