import { GetServerSideProps } from "next";
import PersonPage from "../../../components/PersonPage";
import fetchAPI from "../../../lib/fetch";
import { PartialMember, User } from "../../../lib/types";

interface Props {
  user: User;
  partialMembers: PartialMember[];
}

export default function Index({ user, partialMembers }: Props) {
  return <PersonPage person={user} />;
}

export const getServerSideProps: GetServerSideProps = async (context) => {
  const name = context.params!.user;
  try {
    const [userResponse, partialMembersResponse] = await Promise.allSettled([
      fetchAPI<User>(`/users/${name}`),
      fetchAPI<PartialMember[]>(`/users/${name}/members`),
    ]);
    if (userResponse.status === "rejected")
      throw new Error("Could not fetch user");
    const user = userResponse.value;
    const partialMembers =
      partialMembersResponse.status === "fulfilled"
        ? partialMembersResponse.value
        : [];
    return { props: { user, partialMembers } };
  } catch (e) {
    console.log(e);
    return { notFound: true };
  }
};
