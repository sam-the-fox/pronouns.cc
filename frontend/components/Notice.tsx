import { ReactNode } from "react";
import { ButtonStyle } from "./Button";

export type NoticeStyle = ButtonStyle;

export interface Props {
  header?: string;
  style?: NoticeStyle;
  children?: ReactNode;
}

export default function Notice(props: Props) {
  if (props.style === undefined) {
    return PrimaryNotice(props);
  }

  switch (props.style) {
    case ButtonStyle.primary:
      return PrimaryNotice(props);
    case ButtonStyle.success:
      return SuccessNotice(props);
    case ButtonStyle.danger:
      return DangerNotice(props);
  }
}

function PrimaryNotice(props: Props) {
  return (
    <div className="bg-blue-500 p-4 rounded-md border-blue-600 text-white">
      {props.children}
    </div>
  );
}

function SuccessNotice(props: Props) {
  return (
    <div className="bg-green-600 dark:bg-green-700 p-4 rounded-md text-white border-green-700 dark:border-green-800">
      {props.children}
    </div>
  );
}

function DangerNotice(props: Props) {
  return (
    <div className="bg-red-600 dark:bg-red-700 rounded-md text-red-50 border-red-700 dark:border-red-800 max-w-lg">
      {props.header && (
        <h3 className="uppercase text-sm font-bold border-b border-red-700 dark:border-red-800 px-4 py-3">
          {props.header}
        </h3>
      )}
      <div className={props.header ? "px-4 pt-3 pb-4" : "p-4"}>
        {props.children}
      </div>
    </div>
  );
}
