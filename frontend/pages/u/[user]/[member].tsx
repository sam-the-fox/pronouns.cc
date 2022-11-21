import { GetServerSideProps } from "next";
import fetchAPI from "../../../lib/fetch";
import { Member } from "../../../lib/types";
import PersonPage from "../../../components/PersonPage";

interface Props {
  member: Member;
}

export default function MemberPage({ member }: Props) {
  return <PersonPage person={member} />;
}

export const getServerSideProps: GetServerSideProps = async (context) => {
  try {
    const member = await fetchAPI<Member>(
      `/users/${context.params!.user}/members/${context.params!.member}`
    );

    return { props: { member } };
  } catch (e) {
    console.log(e);

    return { notFound: true };
  }
};
