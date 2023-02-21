import React from "react";

import BackgroundIcons from "@assets/background-icons.svg";
import Background from "@assets/background.svg";

export const OnlyofficeBackground: React.FC = () => (
  <>
    <BackgroundIcons className="absolute bottom-0 -z-25 max-w-none select-none" />
    <Background className="absolute bottom-0 -z-50 max-w-none select-none" />
  </>
);
