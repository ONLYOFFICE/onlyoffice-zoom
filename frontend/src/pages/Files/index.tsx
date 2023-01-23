import React, { useState, useRef, useCallback } from "react";
import { useNavigate } from "react-router-dom";
import { motion } from "framer-motion";

import { OnlyofficeSearchBar } from "@components/search";
import { OnlyofficeTitle } from "@components/title";
import { OnlyofficeFile } from "@components/file";
import { OnlyofficeButton } from "@components/button";
import { OnlyofficeSpinner } from "@components/spinner";
import { OnlyofficeNoFile } from "@components/nofile";

import { useFileSearch } from "@hooks/useFileSearch";

import { formatBytes, getFileIcon, isFileSupported } from "@utils/file";

import { File } from "src/types/file";

export const FilesPage: React.FC = () => {
  const navigate = useNavigate();
  const [query, setQuery] = useState("");
  const {
    isLoading,
    error,
    fetchNextPage,
    isFetchingNextPage,
    files,
    hasNextPage,
  } = useFileSearch(query);
  const observer = useRef<IntersectionObserver>();
  const lastItem = useCallback(
    (node: Element | null) => {
      if (isLoading) return;
      if (observer.current) observer.current.disconnect();
      observer.current = new IntersectionObserver(async (entries) => {
        if (entries[0].isIntersecting && hasNextPage) {
          fetchNextPage();
        }
      });
      if (node) observer.current.observe(node);
    },
    [isLoading, fetchNextPage, hasNextPage]
  );

  const openFile = (file: File) => {
    if (!isFileSupported(file.file_name)) return;
    navigate(
      `/editor?file=${file.file_id}&name=${
        file.file_name
      }&url=${encodeURIComponent(file.download_url)}`
    );
  };

  return (
    <motion.div
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
      className="relative w-full max-w-[790px] h-full flex flex-col my-0 mx-auto md:py-10 pt-10 pb-0"
    >
      <div className="w-full h-20 flex justify-center items-center px-5 pb-10">
        <OnlyofficeButton
          text="Create with ONLYOFFICE"
          primary
          fullWidth
          onClick={() => navigate("/create")}
        />
      </div>
      <div className="table-shadow pb-10 h-[calc(100%-5rem)]">
        <div className="flex items-center justify-center h-12 mx-5 max-w-full truncate text-ellipsis">
          <OnlyofficeTitle text="My Zoom documents" />
        </div>
        <div className="flex h-12 px-5">
          <OnlyofficeSearchBar
            placeholder="Search"
            value={query}
            onChange={(e) => setQuery(e.target.value)}
          />
        </div>
        <div className="px-5 overflow-scroll h-[calc(100%-3rem-1rem)] md:justify-between no-scrollbar">
          {!isLoading && (!!error || files?.length === 0) && (
            <OnlyofficeNoFile title="Could not find zoom files" />
          )}
          {!error &&
            files &&
            files.length > 0 &&
            files?.map((file, index) => {
              if (files.length === index + 1) {
                return (
                  <div key={file.file_id} ref={lastItem}>
                    <OnlyofficeFile
                      icon={getFileIcon(file.file_name)}
                      name={file.file_name}
                      time={new Date(file.timestamp).toLocaleString()}
                      size={formatBytes(+file.file_size)}
                      onClick={() => openFile(file)}
                      supported={isFileSupported(file.file_name)}
                    />
                  </div>
                );
              }
              return (
                <div key={file.file_id}>
                  <OnlyofficeFile
                    icon={getFileIcon(file.file_name)}
                    name={file.file_name}
                    time={new Date(file.timestamp).toLocaleString()}
                    size={formatBytes(+file.file_size)}
                    onClick={() => openFile(file)}
                    supported={isFileSupported(file.file_name)}
                  />
                </div>
              );
            })}
          {(isLoading || isFetchingNextPage) && (
            <div
              className={`relative w-full ${
                isLoading ? "h-full" : "h-fit"
              } my-5 flex justify-center items-center`}
            >
              <OnlyofficeSpinner />
            </div>
          )}
        </div>
      </div>
    </motion.div>
  );
};
export default FilesPage;
