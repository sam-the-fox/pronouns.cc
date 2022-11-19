import { ReactNode } from "react";
import Link from "next/link";

export interface Props {
  children?: ReactNode | undefined;
  href: string;
  plain?: boolean | undefined; // Do not wrap in <li></li>
}

export default function NavItem(props: Props) {
  const ret = (
    <Link
      href={props.href}
      className="hover:text-sky-500 dark:hover:text-sky-400"
    >
      {props.children}
    </Link>
  );

  if (props.plain) {
    return ret;
  }
  return <li>{ret}</li>;
}
