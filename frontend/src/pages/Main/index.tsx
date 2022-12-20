import React from "react";
import { useNavigate } from "react-router-dom";
import { motion } from "framer-motion";

import { OnlyofficeBackground } from "@pages/Main/Background";
import { OnlyofficeMainHeader } from "@pages/Main/Header";
import { OnlyofficeButton } from "@components/button";

export const MainPage: React.FC = () => {
  const navigate = useNavigate();
  return (
    <motion.div initial={{ opacity: 0 }} animate={{ opacity: 1 }}>
      <div className="pt-5 large:pt-12 small:pt-4 xsmall:pt-2">
        <OnlyofficeMainHeader
          title="Welcome to ONLYOFFICE!"
          subtitle="You may open and create a new document without registration"
        />
      </div>
      <div className="flex justify-center items-center flex-col pt-10 pb-10 lg:mx-96 md:mx-40 sm:mx-20 xs:mx-0">
        <div className="pb-5 w-full px-5 z-50">
          <OnlyofficeButton
            text="Create with ONLYOFFICE"
            primary
            fullWidth
            onClick={() => navigate("/create")}
          />
        </div>
        <div className="w-full px-5 z-50">
          <OnlyofficeButton
            text="My Zoom documents"
            fullWidth
            onClick={() => navigate("/files")}
          />
        </div>
      </div>
      <div className="absolute flex justify-center items-center bottom-[0%] w-screen lg:w-4/6 lg:mx-[16%]">
        <OnlyofficeBackground />
      </div>
    </motion.div>
  );
};

export default MainPage;
