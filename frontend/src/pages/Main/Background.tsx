import React from "react";

import background from "@assets/background.svg";

export const OnlyofficeBackground: React.FC = () => (
  <img
    className="relative -z-50 max-w-none scrollable:hidden select-none"
    src={background}
    alt="background"
  />
);
