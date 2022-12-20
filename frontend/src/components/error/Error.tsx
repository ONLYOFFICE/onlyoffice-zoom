import React from "react";

type ErrorProps = {
  text: string;
};

export const OnlyofficeError: React.FC<ErrorProps> = ({ text }) => (
  <div className="flex flex-col justify-center items-center">
    <span className="font-semibold text-center flex items-center">{text}</span>
  </div>
);
