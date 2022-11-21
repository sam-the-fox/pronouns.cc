import { GetServerSideProps } from "next";
import PersonPage from "../../../components/PersonPage";
import fetchAPI from "../../../lib/fetch";
import { PartialMember, User } from "../../../lib/types";

interface Props {
  user: User;
  partialMembers: PartialMember[];
}

export default function Index({ user }: Props) {
  return <PersonPage person={user} />;
}

export const getServerSideProps: GetServerSideProps = async (context) => {
  const name = context.params!.user;
  try {
    const user = await fetchAPI<User>(`/users/${name}`);
    return { props: { user } };
  } catch (e) {
    console.log(e);
    return { notFound: true };
  }
};
