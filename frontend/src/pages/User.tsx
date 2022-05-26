import React, { useState, useEffect } from "react";
import { Link, useParams } from "react-router-dom";
import { ArrowClockwise } from "react-bootstrap-icons";
import ReactMarkdown from "react-markdown";
import { Helmet } from "react-helmet";

import type { User } from "../lib/types";
import fetchAPI from "../lib/fetch";
import FieldCard from "../lib/FieldCard";
import Card from "../lib/Card";
import { userState } from "../lib/store";
import { useRecoilValue } from "recoil";
import Loading from "../lib/Loading";

function UserPage() {
  const params = useParams();

  const [user, setUser] = useState<User>(null);
  const meUser = useRecoilValue(userState);

  useEffect(() => {
    fetchAPI<User>(`/users/${params.username}`).then((res) => {
      setUser(res);
    });
  }, [params.username]);

  if (user == null) {
    return <Loading />;
  }

  return (
    <>
      <Helmet>
        <title>@{user.username} - pronouns.cc</title>
      </Helmet>
      {meUser && meUser.id === user.id && (
        <div className="lg:w-1/3 mx-auto bg-slate-100 dark:bg-slate-700 shadow rounded-md p-2">
          <span>
            You are currently viewing your{" "}
            <span className="font-bold">public</span> profile.
          </span>
          <br />
          <Link
            to="/edit"
            className="hover:underline text-sky-500 dark:text-sky-400"
          >
            Edit your profile
          </Link>
        </div>
      )}
      <div className="container mx-auto">
        <div className="flex flex-col m-2 p-2 lg:flex-row justify-center lg:justify-start items-center space-y-4 lg:space-y-0 lg:space-x-16 lg:items-start border-b border-slate-200 dark:border-slate-700">
          {user.avatar_url && (
            <img className="max-w-xs rounded-full" src={user.avatar_url} />
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
            {user.links?.length && user.fields?.length && (
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
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4 py-2">
          {user.fields?.map((field, index) => (
            <FieldCard key={index} field={field}></FieldCard>
          ))}
          {user.links?.length && (
            <Card title="Links">
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
            </Card>
          )}
        </div>
      </div>
    </>
  );
}

export default UserPage;
