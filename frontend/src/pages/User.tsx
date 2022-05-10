import React, { useState, useEffect } from "react";
import { useParams } from "react-router-dom";
import { ArrowClockwise } from "react-bootstrap-icons";
import { selectorFamily, useRecoilValue } from "recoil";
import axios from "axios";
import type { APIError, User } from "../lib/types";

const userPageState = selectorFamily({
  key: "userPageState",
  get: (username: string) => async () => {
    const res = await axios.get(`/api/v1/users/${username}`);
    if (res.status !== 200) {
      throw res.data as APIError;
    }

    return res.data as User;
  },
});

function UserPage() {
  const params = useParams();

  const [user, setUser] = useState<User>(null);

  useEffect(() => {
    axios.get(`/api/v1/users/${params.username}`).then((res) => {
      if (res.status !== 200) throw res.data as APIError;

      setUser(res.data as User);
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
      <h1 className="text-xl font-bold">
        {user.username} ({user.id})
      </h1>
    </>
  );
}

export default UserPage;
