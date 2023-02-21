import React from "react";
import cx from "classnames";

type TileProps = {
  text: string;
  icon: React.ReactNode;
  size?: "xs" | "sm";
  selected?: boolean;
  onClick?: React.MouseEventHandler<HTMLDivElement>;
  onKeyDown?: React.KeyboardEventHandler<HTMLDivElement>;
};

export const OnlyofficeTile: React.FC<TileProps> = ({
  text,
  icon,
  size = "xs",
  selected = false,
  onClick,
  onKeyDown,
}) => {
  const card = cx({
    "px-5 py-3.5 rounded-lg transform shadow mb-5 outline-none": true,
    "transition duration-100 ease-linear": true,
    "w-56 h-36": true,
    "max-h-36 flex flex-col justify-center": true,
    "hover:-translate-y-[0.125rem] hover:shadow-xl cursor-pointer": !selected,
    "bg-white": !selected,
    "bg-gray-100 border": selected,
  });

  const spn = cx({
    // "text-sm": size === "sm",
    // "text-xs text-[9px]": size === "xs",
    "text-sm": true,
    "font-semibold text-slate-500": true,
    "overflow-hidden whitespace-nowrap inline-block text-ellipsis": true,
  });

  return (
    <div
      role="button"
      tabIndex={0}
      className={card}
      onClick={onClick}
      onKeyDown={onKeyDown}
    >
      <div className="flex items-center justify-center px-1 py-1">
        <div className="relative flex items-end max-w-fit max-h-fit overflow-hidden select-none">
          {icon}
        </div>
      </div>
      <div className="w-full flex items-center justify-center overflow-hidden">
        <span className={spn}>{text}</span>
      </div>
    </div>
  );
};
