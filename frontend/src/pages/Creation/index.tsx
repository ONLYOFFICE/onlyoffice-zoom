import React, { useEffect, useState } from "react";
import { useSnapshot } from "valtio";
import { motion } from "framer-motion";

import { SessionPage } from "@pages/Session";
import { CreatePage } from "@pages/Creation/Creation";

import { SocketState } from "@context/MainContext";

export const CreationPage: React.FC = () => {
  const [session, setSession] = useState(false);
  const { ready, value } = useSnapshot(SocketState);

  useEffect(() => {
    try {
      const sess = JSON.parse(value);
      if (ready && sess?.in_session) setSession(true);
      if (ready && !sess?.in_session) setSession(false);
    } catch (err) {
      setSession(false);
    }
  }, [ready, value]);

  return (
    <motion.div
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
      transition={{ duration: 0.04 }}
      className="h-full overflow-hidden"
    >
      {session && <SessionPage />}
      {!session && <CreatePage />}
    </motion.div>
  );
};

export default CreationPage;
