import React from "react";

import nofile from "@assets/nofile.svg";

type NoFileProps = {
  title: string;
};

export const OnlyofficeNoFile: React.FC<NoFileProps> = ({ title }) => (
  <div className="h-full w-full flex flex-col justify-center items-center">
    <img src={nofile} alt="no files" />
    <span className="font-sans font-bold text-sm max-w-max text-ellipsis truncate">
      {title}
    </span>
  </div>
);
