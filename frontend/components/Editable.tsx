import {
  EmojiLaughing,
  HandThumbsDown,
  HandThumbsUp,
  Heart,
  People,
  Plus,
  Trash3,
} from "react-bootstrap-icons";

import Card from "./Card";
import TextInput from "./TextInput";
import Button, { ButtonStyle } from "./Button";
import { useState } from "react";

export interface EditField {
  id: number;
  name: string;
  pronouns: Record<string, PronounChoice>;
}

export enum PronounChoice {
  favourite,
  okay,
  jokingly,
  friendsOnly,
  avoid,
}

type EditableCardProps = {
  field: EditField;
  onChangeName: React.ChangeEventHandler<HTMLInputElement>;
  onChangePronoun: React.ChangeEventHandler<HTMLInputElement>;
  onAddPronoun(pronoun: string): void;
  onDeletePronoun(e: React.MouseEvent<HTMLButtonElement>, entry: string): void;
  onChangeFavourite(
    e: React.MouseEvent<HTMLButtonElement>,
    entry: string
  ): void;
  onChangeOkay(e: React.MouseEvent<HTMLButtonElement>, entry: string): void;
  onChangeJokingly(e: React.MouseEvent<HTMLButtonElement>, entry: string): void;
  onChangeFriends(e: React.MouseEvent<HTMLButtonElement>, entry: string): void;
  onChangeAvoid(e: React.MouseEvent<HTMLButtonElement>, entry: string): void;
  onClickDelete: React.MouseEventHandler<HTMLButtonElement>;
};

export function EditableCard(props: EditableCardProps) {
  const [input, setInput] = useState("");

  const footer = (
    <div className="flex justify-between">
      <TextInput value={props.field.name} onChange={props.onChangeName} />
      <Button style={ButtonStyle.danger} onClick={props.onClickDelete}>
        <Trash3 aria-hidden className="inline" /> Delete
      </Button>
    </div>
  );

  return (
    <Card title={props.field.name} draggable footer={footer}>
      <ul>
        {Object.keys(props.field.pronouns).map((pronoun, index) => {
          const choice = props.field.pronouns[pronoun];
          return (
            <li className="flex justify-between my-1 items-center" key={index}>
              <TextInput
                value={pronoun}
                prevValue={pronoun}
                onChange={props.onChangePronoun}
              />
              <div className="rounded-md">
                <button
                  type="button"
                  onClick={(e) => props.onChangeFavourite(e, pronoun)}
                  className={`${
                    choice == PronounChoice.favourite
                      ? "bg-slate-500"
                      : "bg-slate-600"
                  } hover:bg-slate-400 p-2`}
                >
                  <Heart />
                </button>
                <button
                  type="button"
                  onClick={(e) => props.onChangeOkay(e, pronoun)}
                  className={`${
                    choice == PronounChoice.okay
                      ? "bg-slate-500"
                      : "bg-slate-600"
                  } hover:bg-slate-400 p-2`}
                >
                  <HandThumbsUp />
                </button>
                <button
                  type="button"
                  onClick={(e) => props.onChangeJokingly(e, pronoun)}
                  className={`${
                    choice == PronounChoice.jokingly
                      ? "bg-slate-500"
                      : "bg-slate-600"
                  } hover:bg-slate-400 p-2`}
                >
                  <EmojiLaughing />
                </button>
                <button
                  type="button"
                  onClick={(e) => props.onChangeFriends(e, pronoun)}
                  className={`${
                    choice == PronounChoice.friendsOnly
                      ? "bg-slate-500"
                      : "bg-slate-600"
                  } hover:bg-slate-400 p-2`}
                >
                  <People />
                </button>
                <button
                  type="button"
                  onClick={(e) => props.onChangeAvoid(e, pronoun)}
                  className={`${
                    choice == PronounChoice.avoid
                      ? "bg-slate-500"
                      : "bg-slate-600"
                  } hover:bg-slate-400 p-2`}
                >
                  <HandThumbsDown />
                </button>
                <button
                  type="button"
                  onClick={(e) => props.onDeletePronoun(e, pronoun)}
                  className="bg-red-600 dark:bg-red-700 hover:bg-red-700 hover:dark:bg-red-800 p-2"
                >
                  <Trash3 />
                </button>
              </div>
            </li>
          );
        })}
        <li className="flex justify-between my-1 items-center">
          <TextInput value={input} onChange={(e) => setInput(e.target.value)} />
          <Button
            style={ButtonStyle.success}
            onClick={() => {
              if (!input || input === "") return;

              props.onAddPronoun(input);
              setInput("");
            }}
          >
            <Plus aria-hidden className="inline" /> Add
          </Button>
        </li>
      </ul>
    </Card>
  );
}
