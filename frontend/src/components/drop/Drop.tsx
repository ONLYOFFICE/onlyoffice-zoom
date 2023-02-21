/* eslint-disable react/jsx-props-no-spreading */
import React, { useState } from "react";
import { useDropzone, DropEvent, FileRejection } from "react-dropzone";
import cx from "classnames";

import Upload from "@assets/upload.svg";

type DragDropProps = {
  onDrop: <T extends File>(
    acceptedFiles: T[],
    fileRejections: FileRejection[],
    event: DropEvent
  ) => Promise<void>;
  errorText?: string;
  uploadingText?: string;
  selectText?: string;
  dragdropText?: string;
  subtext?: string;
  errorTimeout?: number;
};

export const OnlyofficeDragDrop: React.FC<DragDropProps> = ({
  onDrop,
  errorText = "Could not upload your file. Please contact ONLYOFFICE support.",
  uploadingText = "Uploading...",
  selectText = "Select a file",
  dragdropText = "or drag and drop here",
  subtext = "File size is limited",
  errorTimeout = 4000,
}) => {
  const [uploading, setUploading] = useState<boolean>(() => false);
  const [error, setError] = useState<boolean>(false);
  const uploadRef = React.useRef<HTMLInputElement | null>(null);

  const uploadFile = async (
    file: File | undefined,
    event: DropEvent,
    rejection?: FileRejection
  ) => {
    setError(false);
    setUploading(true);
    if (file) {
      try {
        await onDrop([file], rejection ? [rejection] : [], event);
      } catch {
        setError(true);
        setTimeout(() => setError(false), errorTimeout);
      } finally {
        setUploading(false);
      }
    }
  };

  const { getRootProps, getInputProps, isDragActive } = useDropzone({
    onDrop: (files, rejections, event) => {
      uploadFile(files[0], event, rejections[0]);
    },
    noClick: true,
    noKeyboard: true,
  });

  const style = cx({
    "flex flex-col items-center justify-center p-5": true,
    "border-2 border-slate-300 border-dashed rounded-lg": true,
    "bg-transparent bg-opacity-20 text-black": true,
    "transition-all transition-timing-function: ease-in-out": true,
    "transition-duration: 300ms": true,
    "bg-sky-100": isDragActive,
    "bg-emerald-100": uploading,
    "bg-red-100": error,
  });

  return (
    <div className={`${style} w-full h-full`} {...getRootProps()}>
      <Upload />
      {error && (
        <span className="font-sans font-semibold text-sm text-center">
          {errorText}
        </span>
      )}
      {uploading && !error && (
        <span className="font-sans font-semibold text-sm text-center">
          {uploadingText}
        </span>
      )}
      {!uploading && !error && (
        <>
          <input {...getInputProps()} />
          <input
            type="file"
            id="file"
            ref={uploadRef}
            style={{ display: "none" }}
            onChange={(e) => uploadFile(e.target?.files?.[0], e)}
          />
          <div className="font-sans font-semibold text-sm flex flex-wrap justify-center w-full">
            <button
              type="button"
              className="cursor-pointer outline-none border-b-2 border-dashed border-blue-500 text-blue-500 mr-1 max-w-max text-ellipsis truncate"
              onClick={() => uploadRef.current?.click()}
            >
              {selectText}
            </button>
            <span className="text-center max-w-max text-ellipsis truncate">
              {dragdropText}
            </span>
          </div>
          <span className="font-sans font-normal text-xs text-gray-400 text-center max-w-max text-ellipsis truncate">
            {subtext}
          </span>
        </>
      )}
    </div>
  );
};
