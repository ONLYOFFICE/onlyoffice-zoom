import React from "react";
import cx from "classnames";

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
}) => {
  const visible = cx({
    "flex p-2.5 bg-gray-400 rounded-md text-white": true,
    "hover:rounded-xl hover:bg-gray-500 transition-all duration-300": true,
    hidden: !supported,
  });

  return (
    <div className="flex items-center justify-between w-full border-b py-2 my-1">
      <div className="flex items-center justify-start w-2/4">
        <img className="w-[32px] h-[32px]" src={icon} alt={name} />
        <button
          className="text-left font-semibold font-sans md:text-sm text-xs px-2 w-full h-[32px] overflow-hidden text-ellipsis whitespace-nowrap"
          type="button"
          onClick={onClick}
        >
          {name}
        </button>
      </div>
      <div className="flex items-center justify-center w-3/12">
        <span className="overflow-hidden text-ellipsis font-semibold font-sans text-xs text-gray-500">
          {time}
        </span>
      </div>
      <div className="flex items-center justify-center w-2/12">
        <span className="overflow-hidden text-ellipsis font-semibold font-sans text-xs text-gray-500">
          {size}
        </span>
      </div>
      {/* <div className="flex items-center justify-end pl-5">
        <div>
          <button type="button" className={visible} onClick={onClick}>
            <svg
              xmlns="http://www.w3.org/2000/svg"
              className="md:w-4 md:h-4 w-3 h-3"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
              strokeWidth="2"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"
              />
            </svg>
          </button>
        </div>
      </div> */}
    </div>
  );
};
