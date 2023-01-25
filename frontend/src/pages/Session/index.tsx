import React from "react";
import { useNavigate } from "react-router-dom";

import { OnlyofficeBackground } from "@components/background";
import { OnlyofficeMainHeader } from "@components/header";
import { OnlyofficeButton } from "@components/button";

export const SessionPage: React.FC = () => {
  const navigate = useNavigate();
  return (
    <>
      <div className="flex flex-col justify-center items-center h-3/4 px-8">
        <div className="max-w-[363px]">
          <OnlyofficeMainHeader
            title="Welcome to ONLYOFFICE!"
            subtitle="File editing session has been started by another user. Join this session?"
          />
        </div>
        <div className="w-full pt-5 pl-2 z-50 max-w-[544px] flex">
          <div className="w-full flex items-stretch">
            <OnlyofficeButton
              text="Join"
              primary
              fullWidth
              onClick={() => navigate("/editor")}
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

export default SessionPage;
