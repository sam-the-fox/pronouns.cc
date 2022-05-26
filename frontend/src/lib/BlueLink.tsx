import { Link } from "react-router-dom";

export type Props = {
  to: string;
  children?: React.ReactNode;
};

export default function BlueLink({ to, children }: Props) {
  return (
    <Link to={to} className="hover:underline text-sky-500 dark:text-sky-400">
      {children}
    </Link>
  );
}
