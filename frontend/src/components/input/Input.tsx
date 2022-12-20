/* eslint-disable jsx-a11y/label-has-associated-control */
import React from "react";
import cx from "classnames";

type InputProps = {
  text: string;
  value?: string;
  placeholder?: string;
  errorText?: string;
  valid?: boolean;
  textSize?: "sm" | "xs";
  labelSize?: "sm" | "xs";
  autocomplete?: boolean;
  onChange?: React.ChangeEventHandler<HTMLInputElement>;
};

export const OnlyofficeInput: React.FC<InputProps> = ({
  text,
  value,
  placeholder,
  errorText = "Please fill out this field",
  valid = true,
  textSize = "sm",
  labelSize = "xs",
  autocomplete = false,
  onChange,
}) => {
  const istyle = cx({
    "font-normal text-sm text-gray-700 appearance-none block": true,
    "text-xs": textSize === "xs",
    "w-full bg-slate-100 border rounded-sm h-10 px-4": true,
    "border-gray-light": valid,
    "border-red-600": !valid,
  });

  const pstyle = cx({
    hidden: valid,
  });

  return (
    <div>
      <label className={`font-semibold text-${labelSize} text-gray-700 py-2`}>
        {text}
        <span title="required" className="text-red-600">
          {" *"}
        </span>
      </label>
      <input
        value={value}
        placeholder={placeholder}
        className={istyle}
        required
        autoCorrect={autocomplete ? undefined : "off"}
        autoComplete={autocomplete ? undefined : "off"}
        type="text"
        onChange={onChange}
      />
      <p className={`text-red-600 text-xs ${pstyle}`}>{errorText}</p>
    </div>
  );
};
