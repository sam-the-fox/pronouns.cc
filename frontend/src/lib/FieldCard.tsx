import {
  HeartFill,
  HandThumbsUp,
  HandThumbsDown,
  People,
  EmojiLaughing,
} from "react-bootstrap-icons";

import type { Field } from "./types";

export default function FieldCard({ field }: { field: Field }) {
  return (
    <div className=" bg-slate-100 dark:bg-slate-700 rounded-md shadow">
      <h1 className="text-2xl p-2 border-b border-zinc-200 dark:border-slate-800">{field.name}</h1>
      <div className="flex flex-col p-2">
        {field.favourite.map((entry) => (
          <p className="text-lg font-bold">
            <HeartFill className="inline" /> {entry}
          </p>
        ))}
        {field.okay.length !== 0 && (
          <p>
            <HandThumbsUp className="inline" /> {field.okay.join(", ")}
          </p>
        )}
        {field.jokingly.length !== 0 && (
          <p>
            <EmojiLaughing className="inline" /> {field.jokingly.join(", ")}
          </p>
        )}
        {field.friends_only.length !== 0 && (
          <p>
            <People className="inline" /> {field.friends_only.join(", ")}
          </p>
        )}
        {field.avoid.length !== 0 && (
          <p className="text-slate-600 dark:text-slate-400">
            <HandThumbsDown className="inline" /> {field.avoid.join(", ")}
          </p>
        )}
      </div>
    </div>
  );
}
