import React, { useState, useEffect } from "react";
import { useParams } from "react-router-dom";
import { ArrowClockwise } from "react-bootstrap-icons";
import ReactMarkdown from "react-markdown";
import { Helmet } from "react-helmet";

import type { APIError, User } from "../lib/types";
import fetchAPI from "../lib/fetch";
import FieldCard from "../lib/FieldCard";

function UserPage() {
  const params = useParams();

  const [user, setUser] = useState<User>(null);

  useEffect(() => {
    fetchAPI<User>(`/users/${params.username}`).then((res) => {
      setUser(res);
    });
  }, []);

  if (user == null) {
    return (
      <>
        <ArrowClockwise />
        <span>Loading...</span>
      </>
    );
  }

  return (
    <>
      <Helmet>
        <title>@{user.username} - pronouns.cc</title>
      </Helmet>
      <div className="container mx-auto">
        <div className="flex flex-col lg:flex-row justify-center lg:justify-start items-center space-y-4 lg:space-y-0 lg:space-x-16 lg:items-start">
          {user.avatar_url && (
            <img
              className="max-w-max lg:max-w-lg rounded-full"
              src={user.avatar_url}
            />
          )}
          <div className="flex flex-col lg:mx-auto">
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
            {user.links.length !== 0 && (
              <div className="flex flex-col mx-auto lg:ml-auto">
                {user.links.map((link) => (
                  <a
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
          {user.fields.map((field) => (
            <FieldCard field={field}></FieldCard>
          ))}
        </div>
      </div>
    </>
  );
}

export default UserPage;
