import { MouseEventHandler, ReactNode } from "react";

export enum ButtonStyle {
  primary,
  success,
  danger,
}

export interface Props {
  onClick?: MouseEventHandler<HTMLButtonElement>;
  style?: ButtonStyle;
  bold?: boolean;
  disabled?: boolean;
  noRound?: boolean;
  children?: ReactNode;
}

export default function Button(props: Props) {
  if (props.style === undefined) {
    return PrimaryButton(props);
  }

  switch (props.style) {
    case ButtonStyle.primary:
      return PrimaryButton(props);
    case ButtonStyle.success:
      return SuccessButton(props);
    case ButtonStyle.danger:
      return DangerButton(props);
  }
}

function PrimaryButton(props: Props) {
  return (
    <button
      type="button"
      disabled={props.disabled}
      onClick={props.onClick}
      className={`bg-blue-500 dark:bg-blue-500 hover:bg-blue-700 hover:dark:bg-blue-800 p-2 ${
        !props.noRound && "rounded-md"
      } ${props.disabled && "cursor-not-allowed"} text-white`}
    >
      <span className={props.bold ? "font-bold" : ""}>{props.children}</span>
    </button>
  );
}

function SuccessButton(props: Props) {
  return (
    <button
      type="button"
      disabled={props.disabled}
      onClick={props.onClick}
      className={`bg-green-600 dark:bg-green-700 hover:bg-green-700 hover:dark:bg-green-800 p-2 ${
        !props.noRound && "rounded-md"
      } ${props.disabled && "cursor-not-allowed"} text-white`}
    >
      <span className={props.bold ? "font-bold" : ""}>{props.children}</span>
    </button>
  );
}

function DangerButton(props: Props) {
  return (
    <button
      type="button"
      disabled={props.disabled}
      onClick={props.onClick}
      className={`bg-red-600 dark:bg-red-700 hover:bg-red-700 hover:dark:bg-red-800 p-2  ${
        !props.noRound && "rounded-md"
      } ${props.disabled && "cursor-not-allowed"} text-white`}
    >
      <span className={props.bold ? "font-bold" : ""}>{props.children}</span>
    </button>
  );
}
