import React, { useState, useEffect } from "react";
import { motion } from "framer-motion";

import { FilesPage } from "@pages/Files";
import { SessionPage } from "@pages/Session";

import { useWebsocket } from "@context/WebsocketContext";
import { useZoomLanguage } from "@context/LangContext";
import { OnlyofficeSpinner } from "@components/spinner";

export const MainPage: React.FC = () => {
  const [session, setSession] = useState(false);
  const { ready, error, value } = useWebsocket();
  const { loading } = useZoomLanguage();
  useEffect(() => {
    try {
      const sess = JSON.parse(value);
      if (ready && sess?.in_session) setSession(true);
      if (ready && !sess?.in_session) setSession(false);
    } catch (err) {
      setSession(false);
    }
  }, [ready, error, value]);

  return (
    <motion.div
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
      transition={{ duration: 0.04 }}
      className="h-full overflow-hidden"
    >
      {loading && (
        <div className="h-full flex justify-center items-center">
          <OnlyofficeSpinner />
        </div>
      )}
      {session && !loading && <SessionPage />}
      {!session && !loading && <FilesPage />}
    </motion.div>
  );
};

export default MainPage;
