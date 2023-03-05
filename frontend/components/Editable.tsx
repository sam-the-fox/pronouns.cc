import {
  EmojiLaughing,
  HandThumbsDown,
  HandThumbsUp,
  Heart,
  People,
  Plus,
  ThreeDotsVertical,
  Trash3,
} from "react-bootstrap-icons";

import Card from "./Card";
import TextInput from "./TextInput";
import Button, { ButtonStyle } from "./Button";
import { useState } from "react";
import { WordStatus } from "../lib/api-fetch";
import { ReactSortable } from "react-sortablejs";

export interface EditField {
  id: number;
  name: string;
  values: EditFieldValue[];
}

export interface EditFieldValue {
  id: number;
  value: string;
  status: WordStatus;
}

type EditableCardProps = {
  field: EditField;
  onChangeName: React.ChangeEventHandler<HTMLInputElement>;
  onChangePronoun: React.ChangeEventHandler<HTMLInputElement>;
  onAddPronoun(pronoun: string): void;
  onDeletePronoun(e: React.MouseEvent<HTMLButtonElement>, index: number): void;
  onChangeFavourite(
    e: React.MouseEvent<HTMLButtonElement>,
    index: number
  ): void;
  onChangeOkay(e: React.MouseEvent<HTMLButtonElement>, index: number): void;
  onChangeJokingly(e: React.MouseEvent<HTMLButtonElement>, index: number): void;
  onChangeFriends(e: React.MouseEvent<HTMLButtonElement>, index: number): void;
  onChangeAvoid(e: React.MouseEvent<HTMLButtonElement>, index: number): void;
  onClickDelete: React.MouseEventHandler<HTMLButtonElement>;
  onChangeOrder(newState: EditFieldValue[]): void;
};

export function EditableCard(props: EditableCardProps) {
  const [input, setInput] = useState("");

  const footer = (
    <div className="flex flex-col justify-between space-y-2">
      <div className="flex justify-between items-center px-2 pt-2">
        <TextInput value={input} onChange={(e) => setInput(e.target.value)} />
        <Button
          disabled={!input || input === ""}
          style={ButtonStyle.success}
          onClick={() => {
            if (!input || input === "") return;

            props.onAddPronoun(input);
            setInput("");
          }}
        >
          <Plus aria-hidden className="inline" /> Add entry
        </Button>
      </div>
      <div className="flex justify-between p-2 border-t border-zinc-200 dark:border-slate-800">
        <TextInput value={props.field.name} onChange={props.onChangeName} />
        <Button style={ButtonStyle.danger} onClick={props.onClickDelete}>
          <Trash3 aria-hidden className="inline" /> Delete field
        </Button>
      </div>
    </div>
  );

  return (
    <Card title={props.field.name} draggable footer={footer}>
      <ReactSortable
        handle=".entry-handle"
        list={props.field.values}
        setList={props.onChangeOrder}
      >
        {props.field.values.map((value, index) => {
          return (
            <li className="flex justify-between my-1 items-center" key={index}>
              <ThreeDotsVertical className="entry-handle hover:cursor-grab" />
              <TextInput
                value={value.value}
                prevValue={value.value}
                onChange={props.onChangePronoun}
              />
              <div>
                <button
                  type="button"
                  onClick={(e) => props.onChangeFavourite(e, index)}
                  className={`${
                    value.status == WordStatus.Favourite
                      ? "bg-slate-500"
                      : "bg-slate-600"
                  } text-white hover:bg-slate-400 p-2`}
                >
                  <Heart />
                </button>
                <button
                  type="button"
                  onClick={(e) => props.onChangeOkay(e, index)}
                  className={`${
                    value.status == WordStatus.Okay
                      ? "bg-slate-500"
                      : "bg-slate-600"
                  } text-white hover:bg-slate-400 p-2`}
                >
                  <HandThumbsUp />
                </button>
                <button
                  type="button"
                  onClick={(e) => props.onChangeJokingly(e, index)}
                  className={`${
                    value.status == WordStatus.Jokingly
                      ? "bg-slate-500"
                      : "bg-slate-600"
                  } text-white hover:bg-slate-400 p-2`}
                >
                  <EmojiLaughing />
                </button>
                <button
                  type="button"
                  onClick={(e) => props.onChangeFriends(e, index)}
                  className={`${
                    value.status == WordStatus.FriendsOnly
                      ? "bg-slate-500"
                      : "bg-slate-600"
                  } text-white hover:bg-slate-400 p-2`}
                >
                  <People />
                </button>
                <button
                  type="button"
                  onClick={(e) => props.onChangeAvoid(e, index)}
                  className={`${
                    value.status == WordStatus.Avoid
                      ? "bg-slate-500"
                      : "bg-slate-600"
                  } text-white hover:bg-slate-400 p-2`}
                >
                  <HandThumbsDown />
                </button>
                <button
                  type="button"
                  onClick={(e) => props.onDeletePronoun(e, index)}
                  className="bg-red-600 dark:bg-red-700 hover:bg-red-700 hover:dark:bg-red-800 p-2"
                >
                  <Trash3 />
                </button>
              </div>
            </li>
          );
        })}
      </ReactSortable>
    </Card>
  );
}
