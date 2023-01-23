import React, { useEffect, useState } from "react";
import { OnlyofficeSpinner } from "@components/spinner";

import ZoomApp from "./ZoomApp";
import BrowserApp from "./BrowserApp";

function App() {
  const [zoom, setZoom] = useState(false);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    import("@zoom/appssdk").then((zoomSdk) => {
      setTimeout(() => {
        zoomSdk.default
          .config({
            popoutSize: { width: 480, height: 360 },
            capabilities: [
              "expandApp",
              "getAppContext",
              "openUrl",
              "getRunningContext",
            ],
          })
          .then(() => setZoom(true))
          .catch(() => setZoom(false))
          .finally(() => setLoading(false));
      }, 1500);
    });
  }, []);

  return (
    <div className="w-full h-full flex justify-center items-center">
      {loading && <OnlyofficeSpinner />}
      {!loading && zoom && <ZoomApp />}
      {!loading && !zoom && <BrowserApp />}
    </div>
  );
}

export default App;
