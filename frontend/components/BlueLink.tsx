import Link from "next/link";

export type Props = {
  to: string;
  children?: React.ReactNode;
};

export default function BlueLink({ to, children }: Props) {
  return (
    <Link href={to}>
      <a className="hover:underline text-sky-500 dark:text-sky-400">
        {children}
      </a>
    </Link>
  );
}
