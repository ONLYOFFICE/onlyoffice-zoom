import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import { useTranslation } from "react-i18next";

import { OnlyofficeTitle } from "@components/title";
import { OnlyofficeInput } from "@components/input";
import { OnlyofficeTile } from "@components/tile";
import { OnlyofficeButton } from "@components/button";

import { getCreateFileUrl } from "@utils/file";

import Docx from "@assets/docx.svg";
import Pptx from "@assets/pptx.svg";
import Xlsx from "@assets/xlsx.svg";

export const CreatePage: React.FC = () => {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const [file, setFile] = useState(
    t("creation.newfile", "New Document") || "New Document"
  );
  const [fileType, setFileType] = useState<"docx" | "pptx" | "xlsx">("docx");
  const handleChangeFile = (newType: "docx" | "pptx" | "xlsx") => {
    setFileType(newType);
  };

  return (
    <div className="relative w-full max-w-[720px] h-full flex flex-col my-0 mx-auto lg:pt-14 md:pt-7 pt-3">
      <div className="mx-5 md:mx-0 flex h-8 md:h-12 max-w-full truncate text-ellipsis items-center">
        <OnlyofficeTitle text={t("creation.title", "Create with ONLYOFFICE")} />
      </div>
      <div className="px-5 md:px-0 pb-5 h-24">
        <OnlyofficeInput
          text={t("creation.input", "Title")}
          value={file}
          onChange={(e) => setFile(e.target.value)}
          errorText={
            t("creation.input.error") ||
            "File name must contain have 0-200 characters"
          }
          valid={file.length > 0 && file.length <= 200}
        />
      </div>
      <div className="grid md:grid-cols-3 sm:grid-cols-2 xxsmall:grid-cols-1 md:gap-4 sm:gap-2 xsmall:gap-1 justify-center items-center content-start w-full overflow-y-scroll h-[calc(100%-2rem-4rem)] md:h-[calc(100%-3rem-6rem-6rem)] px-5 md:px-0 no-scrollbar">
        <div className="sm:flex sm:justify-center sm:items-center">
          <OnlyofficeTile
            icon={<Docx />}
            text={t("creation.tile.document", "Document")}
            onClick={() => handleChangeFile("docx")}
            onKeyDown={() => handleChangeFile("docx")}
            selected={fileType === "docx"}
          />
        </div>
        <div className="sm:flex sm:justify-center sm:items-center">
          <OnlyofficeTile
            icon={<Xlsx />}
            text={t("creation.tile.spreadsheet", "Spreadsheet")}
            onClick={() => handleChangeFile("xlsx")}
            onKeyDown={() => handleChangeFile("xlsx")}
            selected={fileType === "xlsx"}
          />
        </div>
        <div className="sm:flex sm:justify-center sm:items-center">
          <OnlyofficeTile
            icon={<Pptx />}
            text={t("creation.tile.presentation", "Presentation")}
            onClick={() => handleChangeFile("pptx")}
            onKeyDown={() => handleChangeFile("pptx")}
            selected={fileType === "pptx"}
          />
        </div>
      </div>
      <div className="relative h-16 px-5 md:px-0">
        <OnlyofficeButton
          text={t("button.create", "Create with ONLYOFFICE")}
          disabled={!fileType || !file || file.length > 200}
          fullWidth
          primary
          onClick={() => {
            navigate(
              `/editor?file=${new Date().getTime()}&name=${`${
                encodeURIComponent(file.substring(0, 201)) || "sample"
              }.${fileType}`}&url=${getCreateFileUrl(fileType)}`
            );
          }}
        />
      </div>
      <div className="relative h-16 px-5 md:px-0">
        <OnlyofficeButton
          text={t("button.back", "Back")}
          fullWidth
          onClick={() => {
            navigate(-1);
          }}
        />
      </div>
    </div>
  );
};
