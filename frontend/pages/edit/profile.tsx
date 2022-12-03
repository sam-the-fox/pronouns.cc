import { useRouter } from "next/router";
import { useCallback, useEffect, useState } from "react";
import { useRecoilState } from "recoil";
import Loading from "../../components/Loading";
import cloneDeep from "lodash/cloneDeep";
import { ReactSortable } from "react-sortablejs";

import {
  EditableCard,
  EditField,
  PronounChoice,
} from "../../components/Editable";

import Button, { ButtonStyle } from "../../components/Button";
import { Plus, Save, Trash } from "react-bootstrap-icons";
import toast from "../../lib/toast";

import ReactCodeMirror from "@uiw/react-codemirror";
import { markdown, markdownLanguage } from "@codemirror/lang-markdown";
import ReactMarkdown from "react-markdown";
import { userState } from "../../lib/state";
import { fetchAPI, Field, MeUser } from "../../lib/api-fetch";

export default function Index() {
  const [user, setUser] = useRecoilState(userState);
  const router = useRouter();
  const [state, setState] = useState(cloneDeep(user));

  const originalOrder = state?.fields
    ? state.fields.map((f, i) => {
        const field: EditField = {
          id: i,
          name: f.name,
          pronouns: {},
        };

        f.favourite?.forEach((val) => {
          field.pronouns[val] = PronounChoice.favourite;
        });
        f.okay?.forEach((val) => {
          field.pronouns[val] = PronounChoice.okay;
        });
        f.jokingly?.forEach((val) => {
          field.pronouns[val] = PronounChoice.jokingly;
        });
        f.friends_only?.forEach((val) => {
          field.pronouns[val] = PronounChoice.friendsOnly;
        });
        f.avoid?.forEach((val) => {
          field.pronouns[val] = PronounChoice.avoid;
        });

        return field;
      })
    : [];

  const [fields, setFields] = useState(cloneDeep(originalOrder));

  const resetFields = () => {
    setFields(cloneDeep(originalOrder));
  };

  const addField = () => {
    if (fields.length >= 25) return;

    const lastId = fields[fields.length - 1]?.id ?? 0;

    setFields([...fields, { id: lastId + 1, name: "", pronouns: {} }]);
  };

  useEffect(() => {
    if (!user || !state) {
      router.push("/");
    }
  }, [user]);

  if (!user || !state) {
    return <Loading />;
  }

  const fieldsUpdated = !fieldsEqual(fields, originalOrder);
  const isEdited = fieldsUpdated || state.bio !== user.bio;

  return (
    <div className="container mx-auto">
      <h1 className="p-2 border-b border-slate-300 dark:border-slate-600 flex items-center justify-between">
        <span className="text-3xl">Editing your profile</span>
        {isEdited && (
          <Button
            style={ButtonStyle.success}
            onClick={async () => {
              const user = await updateUser({
                displayName: state.display_name,
                bio: state.bio,
                fields,
              });

              if (user) setUser(user);
            }}
          >
            <Save aria-hidden className="inline" /> Save changes
          </Button>
        )}
      </h1>

      <h3 className="p-2 border-b border-slate-300 dark:border-slate-600 flex items-center justify-between">
        <span className="text-xl">Bio</span>
      </h3>
      <div className="grid grid-cols-1 xl:grid-cols-2 gap-4 py-2">
        <div>
          <h4 className="text-lg font-bold">Edit</h4>
          <ReactCodeMirror
            value={state.bio || undefined}
            onChange={(val, _) => {
              setState({ ...state, bio: val });
            }}
            extensions={[markdown({ base: markdownLanguage })]}
          />
        </div>
        <div>
          <h4 className="text-lg font-bold">Preview</h4>
          <ReactMarkdown>{state.bio || ""}</ReactMarkdown>
        </div>
      </div>

      <h3 className="p-2 border-b border-slate-300 dark:border-slate-600 flex items-center justify-between">
        <span className="text-xl">Fields</span>
        <div className="inline">
          <Button
            noRound
            style={ButtonStyle.success}
            onClick={() => addField()}
          >
            {" "}
            <Plus aria-hidden className="inline" />
            Add field
          </Button>
          {fieldsUpdated && (
            <Button
              noRound
              style={ButtonStyle.danger}
              onClick={() => resetFields()}
            >
              <Trash aria-hidden className="inline" />
              Reset fields
            </Button>
          )}
        </div>
      </h3>
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
            onChangePronoun={(e) => {
              const prev =
                e.target.attributes.getNamedItem("data-prev-value")?.value;
              if (!prev || !e.target.value) return;

              const choice = field.pronouns[prev];
              delete field.pronouns[prev];

              field.pronouns[e.target.value] = choice;

              setFields([...fields]);
            }}
            onAddPronoun={(pronoun) => {
              field.pronouns[pronoun] = PronounChoice.okay;
              setFields([...fields]);
            }}
            onDeletePronoun={(e, pronoun) => {
              delete field.pronouns[pronoun];
              setFields([...fields]);
            }}
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

function fieldsEqual(arr1: EditField[], arr2: EditField[]) {
  if (arr1?.length !== arr2?.length) return false;

  if (!arr1.every((_, i) => arr1[i].id === arr2[i].id)) return false;

  return arr1.every((_, i) =>
    Object.keys(arr1[i].pronouns).every(
      (val) => arr1[i].pronouns[val] === arr2[i].pronouns[val]
    )
  );
}

async function updateUser(args: {
  displayName: string | null;
  bio: string | null;
  fields: EditField[];
}) {
  const newFields = args.fields.map((editField) => {
    const field: Field = {
      name: editField.name,
      favourite: [],
      okay: [],
      jokingly: [],
      friends_only: [],
      avoid: [],
    };

    Object.keys(editField.pronouns).forEach((pronoun) => {
      switch (editField.pronouns[pronoun]) {
        case PronounChoice.favourite:
          field.favourite!.push(pronoun);
          break;
        case PronounChoice.okay:
          field.okay!.push(pronoun);
          break;
        case PronounChoice.jokingly:
          field.jokingly!.push(pronoun);
          break;
        case PronounChoice.friendsOnly:
          field.friends_only!.push(pronoun);
          break;
        case PronounChoice.avoid:
          field.avoid!.push(pronoun);
          break;
      }
    });

    return field;
  });

  try {
    const user = await fetchAPI<MeUser>("/users/@me", "PATCH", {
      display_name: args.displayName ?? null,
      bio: args.bio ?? null,
      fields: newFields,
    });

    toast({ text: "Successfully updated your profile!" });

    return user;
  } catch (e: any) {
    toast({ text: `${e.details ?? e.message ?? e}`, background: "error" });
  }
}
