import React from "react";
import { useNavigate } from "react-router-dom";
import { useTranslation } from "react-i18next";

import { OnlyofficeBackground } from "@components/background";
import { OnlyofficeMainHeader } from "@components/header";
import { OnlyofficeButton } from "@components/button";

export const WelcomePage: React.FC = () => {
  const { t } = useTranslation();
  const navigate = useNavigate();
  return (
    <>
      <div className="flex flex-col justify-center items-center h-3/4 px-8">
        <div className="max-w-[363px] z-50">
          <OnlyofficeMainHeader
            title={t("welcome.title", "Welcome to ONLYOFFICE!")}
            subtitle={t(
              "welcome.subtitle",
              "You may open and create a new document without registration or upload your own files using Drag'n'Drop"
            )}
          />
        </div>
        <div className="w-full pt-5 pl-2 z-50 max-w-[544px] flex">
          <div className="w-full flex items-stretch">
            <OnlyofficeButton
              text={t("button.create", "Create with ONLYOFFICE")}
              primary
              fullWidth
              onClick={() => navigate("/create")}
            />
          </div>
        </div>
      </div>
      <div className="absolute flex justify-center items-center bottom-[0%] w-screen lg:w-4/6 lg:mx-[16%]">
        <OnlyofficeBackground />
      </div>
    </>
  );
};

export default WelcomePage;
