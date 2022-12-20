import React from "react";

import { OnlyofficeTitle, OnlyofficeSubtitle } from "@components/title";

type MainHeaderProps = {
  title: string;
  subtitle: string;
};

export const OnlyofficeMainHeader: React.FC<MainHeaderProps> = ({
  title,
  subtitle,
}) => (
  <div className="flex justify-center items-center flex-col">
    <div className="pb-4">
      <OnlyofficeTitle text={title} large />
    </div>
    <div className="px-5 max-w-[247px]">
      <OnlyofficeSubtitle text={subtitle} />
    </div>
  </div>
);
