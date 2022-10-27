import { GetServerSideProps } from "next";
import fetchAPI from "../../../lib/fetch";
import { Member, User } from "../../../lib/types";
import { userState } from "../../../lib/state";
import { useRecoilValue } from "recoil";

interface Props {
  member: Member;
}

export default function MemberPage({ member }: Props) {
  const isOwnMember = useRecoilValue(userState)?.id === member.user?.id;

  return (
    <>
      <div>hi! this is {isOwnMember ? "" : "not "}your own member.</div>
      <h1>name: {member.name}</h1>
      <p>{member.bio}</p>
    </>
  );
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
