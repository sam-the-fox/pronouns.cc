import {
  HeartFill,
  HandThumbsUp,
  HandThumbsDown,
  People,
  EmojiLaughing,
} from "react-bootstrap-icons";
import BlueLink from "./BlueLink";

import Card from "./Card";
import type { Field } from "./types";

function linkPronoun(input: string) {
  if (input.includes(" ") || input.split("/").length !== 5)
    return <span>{input}</span>;

  const [sub, obj, possDet, possPro, reflexive] = input.split("/");

  return (
    <BlueLink to={`/pronouns/${sub}/${obj}/${possDet}/${possPro}/${reflexive}`}>
      {sub}/{obj}/{possDet}
    </BlueLink>
  );
}

export default function FieldCard({
  field,
  draggable,
}: {
  field: Field;
  draggable?: boolean;
}) {
  return (
    <Card title={field.name} draggable={draggable}>
      {field.favourite.map((entry) => (
        <p className="text-lg font-bold">
          <HeartFill className="inline" /> {linkPronoun(entry)}
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
    </Card>
  );
}
