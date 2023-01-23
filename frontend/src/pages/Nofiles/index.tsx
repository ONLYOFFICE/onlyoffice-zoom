import React from "react";
import { useNavigate } from "react-router-dom";

import { OnlyofficeBackground } from "@pages/Nofiles/Background";
import { OnlyofficeMainHeader } from "@pages/Nofiles/Header";

import { OnlyofficeButton } from "@components/button";

type InitialPageProps = {
  session?: boolean;
};

export const InitialPage: React.FC<InitialPageProps> = ({ session }) => {
  const navigate = useNavigate();
  const subtitle = session
    ? "File editing session has been started by another user. Join this session?"
    : "You may open and create a new document without registration or upload your own files using Drag'n'Drop";
  const btn = session ? "Join" : "Create with ONLYOFFICE";
  return (
    <>
      <div className="flex flex-col justify-center items-center h-3/4 px-8">
        <div className="max-w-[363px]">
          <OnlyofficeMainHeader
            title="Welcome to ONLYOFFICE!"
            subtitle={subtitle}
          />
        </div>
        <div className="w-full pt-5 pl-2 z-50 max-w-[544px] flex">
          <div className="w-full flex items-stretch">
            <OnlyofficeButton
              text={btn}
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

export default InitialPage;
