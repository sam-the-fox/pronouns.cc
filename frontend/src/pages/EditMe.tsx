import cloneDeep from "lodash/cloneDeep";
import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { ReactSortable } from "react-sortablejs";
import { useRecoilValue } from "recoil";

import { userState } from "../lib/store";
import { Field } from "../lib/types";
import FieldCard from "../lib/FieldCard";

interface FieldWithID extends Field {
  id: number;
}

export default function EditMe() {
  const navigate = useNavigate();

  const meUser = useRecoilValue(userState);
  if (!meUser) {
    navigate("/");
    return <>Loading...</>;
  }

  const [state, setState] = useState(cloneDeep(meUser));
  // add an ID to every field (not returned by the API normally, but Sortable needs it)
  const originalOrder = state.fields.map((f, i) => {
    const fID = f as FieldWithID;
    fID.id = i;
    return fID;
  });

  const [fields, setFields] = useState(cloneDeep(originalOrder));
  const fieldsUpdated =
    fields.length !== state.fields.length ||
    !fields.every((_, i) => fields[i].id === originalOrder[i].id);

  return (
    <div className="container mx-auto">
      <div>{`fieldsUpdated: ${fieldsUpdated}`}</div>
      {/* @ts-ignore: This component isn't updated to have a "children" prop yet, but it accepts it just fine. */}
      <ReactSortable
        handle=".handle"
        list={fields}
        setList={setFields}
        className="grid grid-cols-1 md:grid-cols-3 gap-4 py-2"
      >
        {fields.map((field, i) => (
          <FieldCard key={i} field={field} draggable></FieldCard>
        ))}
      </ReactSortable>
    </div>
  );
}
