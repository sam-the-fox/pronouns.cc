import { useRouter } from "next/router";
import { useEffect, useState } from "react";
import { useRecoilState } from "recoil";
import Loading from "../../components/Loading";
import fetchAPI from "../../lib/fetch";
import { userState } from "../../lib/state";
import { MeUser, Field } from "../../lib/types";
import cloneDeep from "lodash/cloneDeep";
import { ReactSortable } from "react-sortablejs";
import Card from "../../components/Card";

import {
  EditableCard,
  EditField,
  PronounChoice,
} from "../../components/Editable";

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

  useEffect(() => {
    if (!user) {
      router.push("/");
    }
  }, [user]);

  if (!user) {
    return <Loading />;
  }

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
  displayName: string;
  bio: string;
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

    Object.keys(editField).forEach((pronoun) => {
      switch (editField.pronouns[pronoun]) {
        case PronounChoice.favourite:
          field.favourite?.push(pronoun);
          break;
        case PronounChoice.okay:
          field.okay?.push(pronoun);
          break;
        case PronounChoice.jokingly:
          field.jokingly?.push(pronoun);
          break;
        case PronounChoice.friendsOnly:
          field.friends_only?.push(pronoun);
          break;
        case PronounChoice.avoid:
          field.avoid?.push(pronoun);
          break;
      }
    });

    return field;
  });

  return await fetchAPI<MeUser>("/users/@me", "PATCH", {
    display_name: args.displayName,
    bio: args.bio,
    fields: newFields,
  });
}
