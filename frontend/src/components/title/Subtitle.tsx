import React from "react";
import cx from "classnames";

type SubtitleProps = {
  text: string;
  large?: boolean;
};

export const OnlyofficeSubtitle: React.FC<SubtitleProps> = ({
  text,
  large = false,
}) => {
  const style = cx({
    "text-slate-800 font-normal text-center": !!text,
    "text-sm": !large,
    "text-base": large,
  });

  return <p className={style}>{text}</p>;
};
