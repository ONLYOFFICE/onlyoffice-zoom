import React from "react";

type FileProps = {
  icon: string;
  name: string;
  time: string;
  size: string;
  supported?: boolean;
  onClick?: React.MouseEventHandler<HTMLButtonElement>;
};

export const OnlyofficeFile: React.FC<FileProps> = ({
  icon,
  name,
  time,
  size,
  supported = false,
  onClick,
}) => (
  <div className="flex items-center justify-between w-full border-b py-2 my-1">
    <div className="flex items-center justify-start w-2/4">
      <img className="w-[32px] h-[32px]" src={icon} alt={name} />
      <button
        className={`${
          supported ? "cursor-pointer" : "cursor-default"
        } text-left font-semibold font-sans md:text-sm text-xs px-2 w-full h-[32px] overflow-hidden text-ellipsis whitespace-nowrap`}
        type="button"
        onClick={onClick}
      >
        {name}
      </button>
    </div>
    <div className="flex items-center justify-center w-3/12">
      <span className="overflow-hidden inline-block font-semibold font-sans text-xs text-gray-500">
        {time}
      </span>
    </div>
    <div className="flex items-center justify-center w-2/12 h-[32px] text-center overflow-hidden">
      <span className="overflow-hidden inline-block font-semibold font-sans text-xs text-gray-500">
        {size}
      </span>
    </div>
  </div>
);
