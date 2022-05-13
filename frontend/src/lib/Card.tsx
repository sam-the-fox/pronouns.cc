import React, { PropsWithChildren } from "react";

export type Props = PropsWithChildren<{ title: string }>;

export default function Card({ title, children }: Props) {
  return (
    <div className="bg-slate-100 dark:bg-slate-700 rounded-md shadow">
      <h1 className="text-2xl p-2 border-b border-zinc-200 dark:border-slate-800">
        {title}
      </h1>
      <div className="flex flex-col p-2">{children}</div>
    </div>
  );
}
