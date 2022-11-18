import { GetServerSideProps } from "next";
import Head from "next/head";
import fetchAPI from "../../../lib/fetch";
import { Name, Pronoun, User, WordStatus } from "../../../lib/types";
import FieldCard from "../../../components/FieldCard";
import Card from "../../../components/Card";
import ReactMarkdown from "react-markdown";
import Image from "next/image";
import { userState } from "../../../lib/state";
import { useRecoilValue } from "recoil";
import Link from "next/link";
import FallbackImage from "../../../components/FallbackImage";
import { ReactNode } from "react";
import {
  EmojiLaughing,
  HandThumbsDown,
  HandThumbsUp,
  HeartFill,
  People,
} from "react-bootstrap-icons";

interface Props {
  user: User;
}

export default function Index({ user }: Props) {
  const isMeUser = useRecoilValue(userState)?.id === user.id;

  return (
    <>
      <Head>
        <title key="title">@{user.username} - pronouns.cc</title>
      </Head>
      {isMeUser && (
        <div className="lg:w-1/3 mx-auto bg-slate-100 dark:bg-slate-700 shadow rounded-md p-2">
          <span>
            You are currently viewing your{" "}
            <span className="font-bold">public</span> profile.
          </span>
          <br />
          <Link
            href="/edit/profile"
            className="hover:underline text-sky-500 dark:text-sky-400"
          >
            Edit your profile
          </Link>
        </div>
      )}
      <div className="container mx-auto">
        <div className="flex flex-col m-2 p-2 lg:flex-row justify-center lg:justify-start items-center space-y-4 lg:space-y-0 lg:space-x-16 lg:items-start border-b border-slate-200 dark:border-slate-700">
          {user.avatar_urls && user.avatar_urls.length !== 0 && (
            <FallbackImage
              className="max-w-xs rounded-full"
              urls={user.avatar_urls}
              alt={`@${user.username}'s avatar`}
            />
          )}
          <div className="flex flex-col">
            {user.display_name && (
              <h1 className="text-2xl font-bold">{user.display_name}</h1>
            )}
            <h3
              className={`${
                user.display_name
                  ? "text-xl italic text-slate-600 dark:text-slate-400"
                  : "text-2xl font-bold"
              }`}
            >
              @{user.username}
            </h3>
            {user.bio && (
              <ReactMarkdown className="prose dark:prose-invert prose-slate">
                {user.bio}
              </ReactMarkdown>
            )}
            {user.links?.length && (
              <div className="flex flex-col mx-auto lg:ml-auto">
                {user.links.map((link, index) => (
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
        </div>
        {user.names?.length > 0 && (
          <div className="border-b border-slate-200 dark:border-slate-700">
            {user.names.map((name, index) => (
              <NameEntry name={name} key={index} />
            ))}
          </div>
        )}
        {user.pronouns?.length > 0 && (
          <div className="border-b border-slate-200 dark:border-slate-700">
            {user.pronouns.map((pronoun, index) => (
              <PronounEntry pronoun={pronoun} key={index} />
            ))}
          </div>
        )}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4 py-2">
          {user.fields?.map((field, index) => (
            <FieldCard key={index} field={field}></FieldCard>
          ))}
        </div>
      </div>
    </>
  );
}

const entryIcon = (status: WordStatus) => {
  let icon: ReactNode;

  switch (status) {
    case WordStatus.Favourite:
      icon = <HeartFill className="inline" />;
      break;
    case WordStatus.Okay:
      icon = <HandThumbsUp className="inline" />;
      break;
    case WordStatus.Jokingly:
      icon = <EmojiLaughing className="inline" />;
      break;
    case WordStatus.FriendsOnly:
      icon = <People className="inline" />;
      break;
    case WordStatus.Avoid:
      icon = <HandThumbsDown className="inline" />;
      break;
  }

  return icon;
};

function NameEntry(props: { name: Name }) {
  const { name } = props;

  return (
    <p
      className={`text-lg ${
        name.status === WordStatus.Favourite && "font-bold"
      }`}
    >
      {entryIcon(name.status)} {name.name}
    </p>
  );
}

function PronounEntry(props: { pronoun: Pronoun }) {
  const { pronoun } = props;

  return (
    <p
      className={`text-lg ${
        pronoun.status === WordStatus.Favourite && "font-bold"
      }`}
    >
      {entryIcon(pronoun.status)}{" "}
      {pronoun.display_text ??
        pronoun.pronouns.split("/").slice(0, 2).join("/")}
    </p>
  );
}

export const getServerSideProps: GetServerSideProps = async (context) => {
  try {
    const user = await fetchAPI<User>(`/users/${context.params!.user}`);

    return { props: { user } };
  } catch (e) {
    console.log(e);

    return { notFound: true };
  }
};
