import { ReactNode, PropsWithChildren } from "react";
import { Link } from "react-router-dom";

export interface Props {
    children?: ReactNode | undefined;
    to: string;
    plain?: boolean | undefined; // Do not wrap in <li></li>
}

export default function NavItem(props: Props) {
    const ret = <Link
        className="hover:text-sky-500 dark:hover:text-sky-400"
        to={props.to}
    >
        {props.children}
    </Link>

    if (props.plain) {
        return ret
    }
    return <li>{ret}</li>;
}
