import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import { motion } from "framer-motion";

import { OnlyofficeTitle } from "@components/title";
import { OnlyofficeInput } from "@components/input";
import { OnlyofficeTile } from "@components/tile";
import { OnlyofficeButton } from "@components/button";

import { getCreateFileUrl } from "@utils/file";

import docx from "@assets/docx.svg";
import pptx from "@assets/pptx.svg";
import xls from "@assets/xls.svg";

export const CreationPage: React.FC = () => {
  const navigate = useNavigate();
  const [file, setFile] = useState("");
  const [fileType, setFileType] = useState<
    "docx" | "pptx" | "xlsx" | undefined
  >(undefined);

  const handleChangeFile = (newType: "docx" | "pptx" | "xlsx") => {
    if (newType === fileType) {
      setFileType(undefined);
      return;
    }
    setFileType(newType);
  };

  return (
    <motion.div
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
      className="relative w-full max-w-[720px] h-full flex flex-col my-0 mx-auto lg:pt-14 md:pt-7 pt-3"
    >
      <div className="mx-5 md:mx-0 flex h-8 md:h-12 max-w-full truncate text-ellipsis items-center">
        <OnlyofficeTitle text="Create with ONLYOFFICE" />
      </div>
      <div className="px-5 md:px-0 pb-5 h-24">
        <OnlyofficeInput
          text="Title"
          value={file}
          onChange={(e) => setFile(e.target.value)}
        />
      </div>
      <div className="grid md:grid-cols-3 sm:grid-cols-2 xxsmall:grid-cols-1 md:gap-4 sm:gap-2 xsmall:gap-1 justify-center items-center content-start w-full overflow-y-scroll h-[calc(100%-2rem-4rem)] md:h-[calc(100%-3rem-6rem-6rem)] px-5 md:px-0">
        <div className="sm:flex sm:justify-center sm:items-center">
          <OnlyofficeTile
            icon={docx}
            text="Document"
            onClick={() => handleChangeFile("docx")}
            onKeyDown={() => handleChangeFile("docx")}
            selected={fileType === "docx"}
          />
        </div>
        <div className="sm:flex sm:justify-center sm:items-center">
          <OnlyofficeTile
            icon={xls}
            text="Spreadsheet"
            onClick={() => handleChangeFile("xlsx")}
            onKeyDown={() => handleChangeFile("xlsx")}
            selected={fileType === "xlsx"}
          />
        </div>
        <div className="sm:flex sm:justify-center sm:items-center">
          <OnlyofficeTile
            icon={pptx}
            text="Presentation"
            onClick={() => handleChangeFile("pptx")}
            onKeyDown={() => handleChangeFile("pptx")}
            selected={fileType === "pptx"}
          />
        </div>
      </div>
      <div className="relative h-16 px-5 md:px-0">
        <OnlyofficeButton
          text="Create with ONLYOFFICE"
          disabled={!fileType || !file}
          fullWidth
          primary
          onClick={() => {
            navigate(
              `/editor?file=${new Date().getTime()}&name=${`${
                encodeURI(file) || "sample"
              }.${fileType}`}&url=${getCreateFileUrl(fileType)}`
            );
          }}
        />
      </div>
    </motion.div>
  );
};

export default CreationPage;
