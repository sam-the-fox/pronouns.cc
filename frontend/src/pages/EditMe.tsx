import cloneDeep from "lodash/cloneDeep";

import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { ReactSortable } from "react-sortablejs";
import { useRecoilValue } from "recoil";
import {
  EmojiLaughing,
  HandThumbsDown,
  HandThumbsUp,
  Heart,
  People,
  Trash3,
} from "react-bootstrap-icons";
import CreatableSelect from "react-select/creatable";

import { userState } from "../lib/store";
import { Field } from "../lib/types";
import Loading from "../lib/Loading";
import Card from "../lib/Card";
import TextInput from "../lib/TextInput";

interface EditField {
  id: number;
  name: string;
  pronouns: Record<string, PronounChoice>;
}

enum PronounChoice {
  favourite,
  okay,
  jokingly,
  friendsOnly,
  avoid,
}

function fieldsEqual(arr1: EditField[], arr2: EditField[]) {
  if (arr1?.length !== arr2?.length) return false;

  if (!arr1.every((_, i) => arr1[i].id === arr2[i].id)) return false;

  return arr1.every((_, i) =>
    Object.keys(arr1[i].pronouns).every(
      (val) => arr1[i].pronouns[val] === arr2[i].pronouns[val]
    )
  );
}

export default function EditMe() {
  const navigate = useNavigate();

  const meUser = useRecoilValue(userState);

  useEffect(() => {
    if (!meUser) {
      navigate("/");
    }
  });
  if (!meUser) {
    return <Loading />;
  }

  const [state, setState] = useState(cloneDeep(meUser));
  // convert all fields to EditFields
  const originalOrder = state.fields.map((f, i) => {
    const field: EditField = {
      id: i,
      name: f.name,
      pronouns: {},
    };

    f.favourite.forEach((val) => {
      field.pronouns[val] = PronounChoice.favourite;
    });
    f.okay.forEach((val) => {
      field.pronouns[val] = PronounChoice.okay;
    });
    f.jokingly.forEach((val) => {
      field.pronouns[val] = PronounChoice.jokingly;
    });
    f.friends_only.forEach((val) => {
      field.pronouns[val] = PronounChoice.friendsOnly;
    });
    f.avoid.forEach((val) => {
      field.pronouns[val] = PronounChoice.avoid;
    });

    return field;
  });

  const [fields, setFields] = useState(cloneDeep(originalOrder));
  const fieldsUpdated = !fieldsEqual(fields, originalOrder);

  return (
    <div className="container mx-auto">
      <div>{`fieldsUpdated: ${fieldsUpdated}`}</div>
      {/* @ts-ignore: This component isn't updated to have a "children" prop yet, but it accepts it just fine. */}
      <ReactSortable
        handle=".handle"
        list={fields}
        setList={setFields}
        className="grid grid-cols-1 xl:grid-cols-2 gap-4 py-2"
      >
        {fields.map((field, i) => (
          <EditableCard
            key={i}
            field={field}
            onChangeName={(e) => {
              field.name = e.target.value;
              setFields([...fields]);
            }}
            onChangeFavourite={(e, entry: string) => {
              field.pronouns[entry] = PronounChoice.favourite;
              setFields([...fields]);
            }}
            onChangeOkay={(e, entry: string) => {
              field.pronouns[entry] = PronounChoice.okay;
              setFields([...fields]);
            }}
            onChangeJokingly={(e, entry: string) => {
              field.pronouns[entry] = PronounChoice.jokingly;
              setFields([...fields]);
            }}
            onChangeFriends={(e, entry: string) => {
              field.pronouns[entry] = PronounChoice.friendsOnly;
              setFields([...fields]);
            }}
            onChangeAvoid={(e, entry: string) => {
              field.pronouns[entry] = PronounChoice.avoid;
              setFields([...fields]);
            }}
            onClickDelete={(_) => {
              const newFields = [...fields];
              newFields.splice(i, 1);
              setFields(newFields);
            }}
          />
        ))}
      </ReactSortable>
    </div>
  );
}

type EditableCardProps = {
  field: EditField;
  onChangeName: React.ChangeEventHandler<HTMLInputElement>;
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

function EditableCard(props: EditableCardProps) {
  const footer = (
    <div className="flex justify-between">
      <TextInput value={props.field.name} onChange={props.onChangeName} />
      <button
        type="button"
        onClick={props.onClickDelete}
        className="bg-red-600 dark:bg-red-700 hover:bg-red-700 hover:dark:bg-red-800 p-2 rounded-md"
      >
        <Trash3 aria-hidden className="inline" />{" "}
        <span className="font-bold">Delete</span>
      </button>
    </div>
  );

  return (
    <Card title={props.field.name} draggable footer={footer}>
      <ul>
        {Object.keys(props.field.pronouns).map((pronoun, index) => {
          const choice = props.field.pronouns[pronoun];
          return (
            <li className="flex justify-between my-1" key={index}>
              <div>{pronoun}</div>
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
                  className="bg-red-600 dark:bg-red-700 hover:bg-red-700 hover:dark:bg-red-800 p-2"
                >
                  <Trash3 />
                </button>
              </div>
            </li>
          );
        })}
      </ul>
    </Card>
  );
}
