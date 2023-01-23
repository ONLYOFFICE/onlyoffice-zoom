import React from "react";

import { OnlyofficeNavbar } from "@layouts/navbar";

type BasicLayoutProps = {
  onNavbarClick?: React.MouseEventHandler<HTMLButtonElement>;
  children?: JSX.Element | JSX.Element[];
};

export const OnlyofficeBasicLayoutContainer: React.FC<BasicLayoutProps> = ({
  onNavbarClick,
  children,
}) => (
  <div className="relative h-[100%] w-screen overflow-hidden">
    <div className="h-12">
      <OnlyofficeNavbar onClick={onNavbarClick} />
    </div>
    <div className="h-[calc(100%-3rem)]">{children}</div>
  </div>
);
