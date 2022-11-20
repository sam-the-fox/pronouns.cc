import { GetServerSideProps } from "next";
import Head from "next/head";
import fetchAPI from "../../../lib/fetch";
import { Field, Name, PartialMember, Pronoun, User, WordStatus } from "../../../lib/types";
import ReactMarkdown from "react-markdown";
import { userState } from "../../../lib/state";
import { useRecoilValue } from "recoil";
import FallbackImage from "../../../components/FallbackImage";
import {
  EmojiLaughing,
  HandThumbsDown,
  HandThumbsUp,
  HeartFill,
  People,
} from "react-bootstrap-icons";
import BlueLink from "../../../components/BlueLink";
import React from "react";
import Card from "../../../components/Card";

interface Props {
  user: User;
  partialMembers: PartialMember[];
}

export default function Index({ user, partialMembers }: Props) {
  return (
    <>
      <Head>
        <title key='title'>{`@${user.username} - pronouns.cc`}</title>
      </Head>
      <IsOwnPageNotice user={user} />
      <div className="container mx-auto pb-[20vh]">
        <div className="
          m-2 p-2
          flex flex-col lg:flex-row
          justify-center lg:justify-start
          items-center lg:items-start
          lg:space-x-16
          space-y-4 lg:space-y-0
          border-b border-slate-200 dark:border-slate-700
        ">
          <UserAvatar user={user} />
          <UserInfo user={user} />
        </div>
        <LabelList source={user.names} />
        <LabelList source={user.pronouns} />
        <FieldCardGrid fields={user.fields} />
        <MemberList user={user} partialMembers={partialMembers} />
      </div>
    </>
  );
}

export const getServerSideProps: GetServerSideProps = async (context) => {
  const username = context.params!.user;
  try {
    const [userResponse, partialMembersResponse] = await Promise.allSettled([
      fetchAPI<User>(`/users/${username}`),
      fetchAPI<PartialMember[]>(`/users/${username}/members`)
    ]);
    if (userResponse.status === 'rejected')
      throw new Error('Could not fetch user');
    const user = userResponse.value;
    const partialMembers = partialMembersResponse.status === 'fulfilled'
      ? partialMembersResponse.value : [];
    return { props: { user, partialMembers } };
  } catch (e) {
    console.log(e);
    return { notFound: true };
  }
};

function IsOwnPageNotice({ user }: { user: User }) {
  const isThisMyPage = useRecoilValue(userState)?.id === user.id;
  return (
    isThisMyPage || true ? (
      <div className="lg:w-1/3 mx-auto bg-slate-100 dark:bg-slate-700 shadow rounded-md p-2">
        You are currently viewing your <b>public</b> profile.
        <br />
        <BlueLink to="/edit/profile">Edit your profile</BlueLink>
      </div>
    ) : <></>
  );
}

function MemberList({
  user,
  partialMembers,
  className
}: {
  user: User,
  partialMembers: PartialMember[],
  className?: string
}) {
  console.log(partialMembers);
  return (
    <div className={`mx-auto flex-col items-center ${className || ''}`}>
      <h1 className='text-2xl'>Members</h1>
      <ul>
        {partialMembers.map(partialMember =>
          <li className='before:[content:"-_"]' key={partialMember.id}>
            <BlueLink to={`/u/${user.username}/${partialMember.name}`}>
              <span>{partialMember.display_name ?? partialMember.name}</span>
            </BlueLink>
          </li>
        )}
      </ul>
    </div>
  );
}

function UserAvatar({ user }: { user: User }) {
  return (
    user.avatar_urls && user.avatar_urls.length !== 0 ? (
      <FallbackImage
        className="max-w-xs rounded-full"
        urls={user.avatar_urls}
        alt={`@${user.username}'s avatar`}
      />
    ) : <></>
  );
}

function UserInfo({ user }: { user: User }) {
  const { display_name, username, bio, links } = user;
  return (
    <div className="flex flex-col">
      {/* display name */}
      {display_name && (
        <h1 className="text-2xl font-bold">{display_name}</h1>
      )}
      {/* username */}
      <h3
        className={`${
          display_name
            ? "text-xl italic text-slate-600 dark:text-slate-400"
            : "text-2xl font-bold"
        }`}
      >
        @{username}
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

function LabelList({ source }: { source: Name[] | Pronoun[] }) {
  return (
    source?.length > 0 ? (
      <div className="border-b border-slate-200 dark:border-slate-700">
        {source.map((label, index) => (
          <LabelLine key={index} label={label} />
        ))}
      </div>
    ) : <></>
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
    { className: 'inline' }
  );
}

function LabelsLine({ labels }: { labels: Name[] | Pronoun[] }) {
  if (labels.length === 0) return <></>;
  const status = labels[0].status;
  const text = labels
    .map(label =>
      'name' in label
        ? label.name
        : label.display_text ?? label.pronouns.split('/').slice(0, 2).join('/'))
    .join(', ');
  return (
    <p className={`
      ${status === WordStatus.Favourite ? 'text-lg font-bold' : ''}
      ${status === WordStatus.Avoid ? 'text-slate-600 dark:text-slate-400' : ''}`}>
      <LabelStatusIcon status={status} /> {text}
    </p>
  );
}

function LabelLine({ label }: { label: Name | Pronoun }) {
  return <LabelsLine labels={[label] as Name[] | Pronoun[]} />;
}

function FieldCardGrid({ fields }: { fields: Field[] }) {
  return (
    <div className="grid grid-cols-1 md:grid-cols-3 gap-4 py-2">
      {fields?.map((field, index) => (
        <FieldCard field={field} key={index} />
      ))}
    </div>
  );
}

const fieldEntryStatus: { [key in string]: WordStatus } = {
  favourite:    WordStatus.Favourite,
  okay:         WordStatus.Okay,
  jokingly:     WordStatus.Jokingly,
  friends_only: WordStatus.FriendsOnly,
  avoid:        WordStatus.Avoid,
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
      {Object.entries(fieldEntryStatus).map(([statusName, status], i) =>
        <LabelsLine key={i} labels={((field as any)[statusName])?.map((name: string) => ({ name, status }))} />
      )}
    </Card>
  );
}
