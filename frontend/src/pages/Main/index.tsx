import React, { useState, useEffect } from "react";
import { motion } from "framer-motion";

import { OnlyofficeSpinner } from "@components/spinner";

import { fetchFiles } from "@services/file";

import { FilesPage } from "@pages/Files";
import { InitialPage } from "@pages/Nofiles";

export const MainPage: React.FC = () => {
  const [initial, setInitial] = useState(true);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchFiles()
      .then((files) => setInitial(files.response.length < 1))
      .finally(() => setLoading(false));
  }, []);

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
      {!loading && initial && <InitialPage />}
      {!loading && !initial && <FilesPage />}
    </motion.div>
  );
};

export default MainPage;
