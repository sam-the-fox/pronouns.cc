import React, { ReactNode } from "react";

export type Props = {
  children?: ReactNode | undefined;
  title: string;
  draggable?: boolean;
  footer?: ReactNode | undefined;
};

export default function Card({ title, draggable, children, footer }: Props) {
  return (
    <div className="bg-slate-100 dark:bg-slate-700 rounded-md shadow">
      <h1
        className={`text-2xl p-2 border-b border-zinc-200 dark:border-slate-800${
          draggable && " handle hover:cursor-grab"
        }`}
      >
        {title}
      </h1>
      <div className="flex flex-col p-2">{children}</div>
      {footer && (
        <div className="p-2 border-t border-zinc-200 dark:border-slate-800">
          {footer}
        </div>
      )}
    </div>
  );
}
