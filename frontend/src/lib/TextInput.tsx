import { ChangeEventHandler } from "react";

export type Props = {
  defaultValue?: string;
  value?: string;
  onChange?: ChangeEventHandler<HTMLInputElement>;
};

export default function TextInput(props: Props) {
  return (
    <input
      type="text"
      className="p-1 lg:p-2 rounded-md bg-white border-slate-300 text-black dark:bg-slate-800 dark:border-slate-900 dark:text-white"
      defaultValue={props.defaultValue}
      value={props.value}
      onChange={props.onChange}
    />
  );
}
