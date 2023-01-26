import React, { useState, useEffect } from "react";
import { motion } from "framer-motion";

import { FilesPage } from "@pages/Files";
import { SessionPage } from "@pages/Session";

import { useWebsocket } from "@context/WebsocketContext";

export const MainPage: React.FC = () => {
  const [session, setSession] = useState(false);
  const { ready, error, value } = useWebsocket();
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
      {session && <SessionPage />}
      {!session && <FilesPage />}
    </motion.div>
  );
};

export default MainPage;
