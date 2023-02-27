import React, { useEffect, useState } from "react";
import { useTranslation } from "react-i18next";

import { OnlyofficeBasicLayoutContainer } from "@layouts/container";

import { OnlyofficeSpinner } from "@components/spinner";
import { OnlyofficeButton } from "@components/button";
import { OnlyofficeError } from "@components/error";
import { OnlyofficeSubtitle } from "@components/title";

import BackgroundError from "@assets/background-error.svg";

import ZoomApp from "./ZoomApp";
import BrowserApp from "./BrowserApp";

function App() {
  const { t } = useTranslation();
  const [zoom, setZoom] = useState(false);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(false);

  useEffect(() => {
    const offlineHandler = () => setError(true);
    window.addEventListener("offline", offlineHandler);
    import("@zoom/appssdk").then((zoomSdk) => {
      setTimeout(() => {
        zoomSdk.default
          .config({
            popoutSize: { width: 480, height: 360 },
            capabilities: ["expandApp", "getAppContext", "getMeetingUUID"],
          })
          .then(() => setZoom(true))
          .catch(() => setZoom(false))
          .finally(() => setLoading(false));
      }, 1500);
    });

    return () => window.removeEventListener("offline", offlineHandler);
  }, []);

  return (
    <div className="w-full h-full flex justify-center items-center">
      {!loading && error && (
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
      )}
      {loading && !error && (
        <OnlyofficeBasicLayoutContainer>
          <div className="w-full h-full flex justify-center items-center">
            <OnlyofficeSpinner />
          </div>
        </OnlyofficeBasicLayoutContainer>
      )}
      {!loading && !error && zoom && <ZoomApp />}
      {!loading && !error && !zoom && <BrowserApp />}
    </div>
  );
}

export default App;
