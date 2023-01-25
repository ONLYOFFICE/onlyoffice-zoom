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
  <div className="flex justify-center items-center flex-col cursor-default">
    <div className="pb-4">
      <OnlyofficeTitle text={title} large />
    </div>
    <div>
      <OnlyofficeSubtitle text={subtitle} />
    </div>
  </div>
);
