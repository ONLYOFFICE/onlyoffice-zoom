import React, { useState, useEffect } from "react";
import { motion } from "framer-motion";

import { FilesPage } from "@pages/Files";
import { WelcomePage } from "@pages/Welcome";
import { SessionPage } from "@pages/Session";

import { OnlyofficeSpinner } from "@components/spinner";

import { fetchFiles } from "@services/file";

import { useWebsocket } from "@context/WebsocketContext";

export const MainPage: React.FC = () => {
  const [session, setSession] = useState(false);
  const [initial, setInitial] = useState(true);
  const [loading, setLoading] = useState(true);
  const { ready, error, value } = useWebsocket();

  useEffect(() => {
    fetchFiles()
      .then((files) => {
        setInitial(files.response.length < 1);
        setLoading(false);
      })
      .catch(() => setLoading(false));
  }, []);

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
        <div className="h-full w-full flex justify-center items-center">
          <OnlyofficeSpinner />
        </div>
      )}
      {!loading && session && <SessionPage />}
      {!loading && !session && initial && <WelcomePage />}
      {!loading && !session && !initial && <FilesPage />}
    </motion.div>
  );
};

export default MainPage;
