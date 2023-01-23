import React from "react";
import cx from "classnames";

type ButtonProps = {
  text: string;
  disabled?: boolean;
  primary?: boolean;
  fullWidth?: boolean;
  onClick?: React.MouseEventHandler<HTMLButtonElement>;
};

export const OnlyofficeButton: React.FC<ButtonProps> = ({
  text,
  disabled = false,
  primary = false,
  fullWidth = false,
  onClick,
}) => {
  const classes = cx({
    "hover:shadow-lg duration-200": !disabled,
    "bg-sky-500 text-slate-200": primary,
    "bg-white text-black border-2 border-slate-300 border-solid": !primary,
    "w-full": fullWidth,
    "bg-opacity-50 cursor-default": disabled,
  });

  return (
    <button
      type="button"
      disabled={disabled}
      className={`p-3 text-xs lg:text-base font-semibold font-sans rounded-md cursor-pointer ${classes} truncate text-ellipsis`}
      onClick={onClick}
    >
      {text}
    </button>
  );
};
