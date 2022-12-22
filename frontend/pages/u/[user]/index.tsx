import { GetServerSideProps } from "next";

import PersonPage from "../../../components/PersonPage";
import { User } from "../../../lib/api";
import * as API from "../../../lib/api-fetch";

interface Props {
  user: API.User;
}

export default function Index({ user }: Props) {
  return <PersonPage person={new User(user)} />;
}

export const getServerSideProps: GetServerSideProps = async (context) => {
  const userName = context.params!.user;
  if (typeof userName !== "string") return { notFound: true };
  try {
    return { props: { user: await API.fetchAPI<User>(`/users/${userName}`) } };
  } catch (e) {
    console.log(e);
    return { notFound: true };
  }
};
