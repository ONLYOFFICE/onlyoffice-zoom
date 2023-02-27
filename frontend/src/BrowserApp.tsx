import React from "react";
import { useTranslation } from "react-i18next";

import { OnlyofficeButton } from "@components/button";
import { OnlyofficeError } from "@components/error";
import { OnlyofficeSubtitle } from "@components/title";

import BackgroundError from "@assets/background-error.svg";

function BrowserApp() {
  const { t } = useTranslation();
  return (
    <div className="w-screen h-screen flex justify-center flex-col items-center mb-1">
      <div className="absolute flex justify-center items-center w-screen h-screen">
        <BackgroundError />
      </div>
      <div className="pb-5">
        <OnlyofficeError text={t("context.error.title") || "Error"} />
      </div>
      <OnlyofficeSubtitle
        text={
          t("context.error.text") ||
          "Something went wrong. Please reload the page or contact the site administrator."
        }
      />
      <div className="pt-5 z-[100]">
        <OnlyofficeButton
          primary
          text={t("button.reload") || "Reload"}
          onClick={() => window.location.reload()}
        />
      </div>
    </div>
  );
}

export default BrowserApp;
