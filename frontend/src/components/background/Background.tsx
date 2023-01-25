import React from "react";

import backgroundIcons from "@assets/background-icons.svg";
import background from "@assets/background.svg";

export const OnlyofficeBackground: React.FC = () => (
  <>
    <img
      className="absolute bottom-0 -z-25 max-w-none scrollable:hidden select-none"
      src={backgroundIcons}
      alt="background-icons"
    />
    <img
      className="absolute bottom-0 -z-50 max-w-none select-none"
      src={background}
      alt="background"
    />
  </>
);
