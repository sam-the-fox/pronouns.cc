import { markdown, markdownLanguage } from "@codemirror/lang-markdown";
import { githubDark, githubLight } from "@uiw/codemirror-theme-github";
import ReactCodeMirror from "@uiw/react-codemirror";
import cloneDeep from "lodash/cloneDeep";
import { useRouter } from "next/router";
import { useEffect, useState } from "react";
import { Plus, Save, Trash } from "react-bootstrap-icons";
import ReactMarkdown from "react-markdown";
import { ReactSortable } from "react-sortablejs";
import { useRecoilState, useRecoilValue } from "recoil";

import Button, { ButtonStyle } from "../../components/Button";
import { EditableCard, EditField } from "../../components/Editable";
import Loading from "../../components/Loading";
import { fetchAPI, Field, MeUser, WordStatus } from "../../lib/api-fetch";
import { themeState, userState } from "../../lib/state";
import toast from "../../lib/toast";

export default function Index() {
  const [user, setUser] = useRecoilState(userState);
  const darkTheme = useRecoilValue(themeState);
  const router = useRouter();
  const [state, setState] = useState(cloneDeep(user));

  const originalOrder = state?.fields
    ? state.fields.map((f, i) => {
        const field: EditField = {
          id: i,
          name: f.name,
          values: [],
        };

        f.entries?.forEach((entry) => {
          field.values.push(entry);
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

    const lastId = fields[fields.length - 1]?.id ?? -1;

    setFields([
      ...fields,
      { id: lastId + 1, name: `Field #${lastId + 2}`, values: [] },
    ]);
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
            className="text-base"
            value={state.bio || undefined}
            onChange={(val, _) => {
              setState({ ...state, bio: val });
            }}
            theme={darkTheme ? githubDark : githubLight}
            minHeight="200"
            basicSetup={{
              lineNumbers: false,
              bracketMatching: false,
              closeBrackets: false,
              autocompletion: false,
              allowMultipleSelections: false,
            }}
            lang="markdown"
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

              const idx = field.values.findIndex((val) => val.value === prev);
              if (idx !== -1) {
                field.values[idx].value = e.target.value;
              }

              setFields([...fields]);
            }}
            onAddPronoun={(pronoun) => {
              field.values.push({ value: pronoun, status: WordStatus.Okay });
              setFields([...fields]);
            }}
            onDeletePronoun={(e, index) => {
              delete field.values[index];
              setFields([...fields]);
            }}
            onChangeName={(e) => {
              field.name = e.target.value;
              setFields([...fields]);
            }}
            onChangeFavourite={(e, index) => {
              field.values[index].status = WordStatus.Favourite;
              setFields([...fields]);
            }}
            onChangeOkay={(e, index) => {
              field.values[index].status = WordStatus.Okay;
              setFields([...fields]);
            }}
            onChangeJokingly={(e, index) => {
              field.values[index].status = WordStatus.Jokingly;
              setFields([...fields]);
            }}
            onChangeFriends={(e, index) => {
              field.values[index].status = WordStatus.FriendsOnly;
              setFields([...fields]);
            }}
            onChangeAvoid={(e, index) => {
              field.values[index].status = WordStatus.Avoid;
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
    arr1[i].values.every(
      (val, j) =>
        val.value === arr2[i].values[j].value &&
        val.status === arr2[i].values[j].status
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
      entries: [],
    };

    field.entries = [...editField.values];

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
