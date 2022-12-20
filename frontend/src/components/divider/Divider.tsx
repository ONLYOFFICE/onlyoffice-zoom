import React from "react";

type DividerProps = {
  text: string;
};

export const OnlyofficeDivider: React.FC<DividerProps> = ({ text }) => (
  <div className="relative flex py-5 items-center w-full">
    <div className="flex-grow border-t border-gray-400" />
    <span className="flex-shrink mx-4 text-gray-400">{text}</span>
    <div className="flex-grow border-t border-gray-400" />
  </div>
);
