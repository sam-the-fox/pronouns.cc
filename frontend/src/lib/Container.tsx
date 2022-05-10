import React from "react";

export default function Container(props: React.PropsWithChildren<{}>) {
  return <div className="m-2 lg:m-4">{props.children}</div>;
}
