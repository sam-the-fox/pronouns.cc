import { ChangeEventHandler } from "react";

export type Props = {
  contrastBackground?: boolean;
  prevValue?: string;
  defaultValue?: string;
  value?: string;
  onChange?: ChangeEventHandler<HTMLInputElement>;
};

export default function TextInput(props: Props) {
  const bg = props.contrastBackground
    ? "bg-slate-50 dark:bg-slate-700"
    : "bg-white dark:bg-slate-800";

  return (
    <input
      type="text"
      data-prev-value={props.prevValue ?? props.value}
      className={`p-1 lg:p-2 rounded-md ${bg} border-slate-300 text-black dark:border-slate-900 dark:text-white`}
      defaultValue={props.defaultValue}
      value={props.value}
      onChange={props.onChange}
    />
  );
}
