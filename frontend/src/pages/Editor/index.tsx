import React from "react";
import { useNavigate, useSearchParams } from "react-router-dom";
import { DocumentEditor } from "@onlyoffice/document-editor-react";
import { motion } from "framer-motion";

import { OnlyofficeSpinner } from "@components/spinner";
import { OnlyofficeError } from "@components/error";
import { OnlyofficeButton } from "@components/button";

import { useBuildConfig } from "@hooks/useBuildConfig";
import { removeSession } from "@services/session";

import icon from "@assets/nofile.svg";

const onEditor = () => {
  const loader = document.getElementById("eloader");
  if (loader) {
    loader.classList.add("opacity-0");
    loader.classList.add("-z-100");
    loader.classList.add("hidden");
  }

  const editor = document.getElementById("editor");
  if (editor) {
    editor.classList.remove("opacity-0");
  }
};

export const OnlyofficeEditorPage: React.FC = () => {
  const [params] = useSearchParams();
  const navigate = useNavigate();
  const { isLoading, error, data } = useBuildConfig(
    params.get("file") || "sample.docx",
    params.get("name") || "sample.docx",
    params.get("url") ||
      "https://d2nlctn12v279m.cloudfront.net/assets/docs/samples/new.docx"
  );

  const validConfig = !error && !isLoading && data;
  return (
    <motion.div
      className="w-screen h-screen"
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
    >
      {!error && (
        <div
          id="eloader"
          className="relative w-screen h-screen flex flex-col small:justify-between justify-center items-center transition duration-250 ease-linear"
        >
          <div className="pb-5 small:h-full small:flex small:items-center">
            <OnlyofficeSpinner />
          </div>
          <div className="small:mb-5 small:px-5 small:w-full">
            <OnlyofficeButton
              primary
              text="Cancel"
              fullWidth
              onClick={() => navigate("/")}
            />
          </div>
        </div>
      )}
      {!!error && (
        <div className="w-screen h-screen flex justify-center flex-col items-center mb-1">
          <img src={icon} alt="error-icon" />
          <OnlyofficeError text="Could not open the file. Something went wrong" />
          <div className="pt-5">
            <OnlyofficeButton
              primary
              text="Go back"
              onClick={() => navigate("/")}
            />
          </div>
        </div>
      )}
      {validConfig && process.env.DOC_SERVER && (
        <div
          id="editor"
          className="w-screen h-screen opacity-0 transition duration-250 ease-linear"
        >
          <DocumentEditor
            id="docxEditor"
            documentServerUrl={process.env.DOC_SERVER}
            config={{
              document: {
                fileType: data.document.fileType,
                key: data.document.key,
                title: data.document.title,
                url: data.document.url,
                permissions: data.document.permissions,
              },
              documentType: data.documentType,
              editorConfig: {
                callbackUrl: data.editorConfig.callbackUrl,
                user: data.editorConfig.user,
                lang: data.editorConfig.lang,
                customization: {
                  goback: {
                    requestClose: true,
                    text: "Close",
                  },
                },
              },
              token: data.token,
              type: data.type,
              events: {
                onRequestClose: () => navigate("/"),
                onAppReady: onEditor,
                onError: () => {
                  onEditor();
                  if (data.is_owner) removeSession();
                },
                onWarning: onEditor,
              },
            }}
          />
        </div>
      )}
    </motion.div>
  );
};

export default OnlyofficeEditorPage;
