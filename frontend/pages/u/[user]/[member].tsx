import { GetServerSideProps } from "next";
import PersonPage from "../../../components/PersonPage";
import { Member } from "../../../lib/api";
import * as API from "../../../lib/api-fetch";

interface Props {
  member: API.Member;
}

export default function MemberPage({ member }: Props) {
  return <PersonPage person={new Member(member)} />;
}

export const getServerSideProps: GetServerSideProps = async (context) => {
  const userName = context.params!.user;
  if (typeof userName !== "string") return { notFound: true };
  const memberName = context.params!.member;
  if (typeof memberName !== "string") return { notFound: true };

  try {
    return {
      props: {
        member: await API.fetchAPI<API.Member>(
          `/users/${context.params!.user}/members/${context.params!.member}`
        ),
      },
    };
  } catch (e) {
    console.log(e);
    return { notFound: true };
  }
};
