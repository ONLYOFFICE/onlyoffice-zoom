import React, { MouseEventHandler } from "react";
import cx from "classnames";

import { OnlyofficeButton } from "@components/button";

type PopupProps = {
  visible: boolean;
  title: string;
  text: string;
  mainBtn: string;
  secBtn: string;
  close: () => void;
  mainAction?: MouseEventHandler;
  secAction?: MouseEventHandler;
};

export const OnlyofficePopup: React.FC<PopupProps> = ({
  visible,
  title,
  text,
  mainBtn,
  secBtn,
  close,
  mainAction,
  secAction,
}) => {
  const cstyles = cx({
    "z-50 pin overflow-auto bg-smoke-light opacity-100": visible,
    "-z-50 opacity-0": !visible,
  });

  const mstyle = cx({
    "opacity-100": visible,
    "opaicty-0": !visible,
  });

  return (
    <div
      className={`fixed top-0 bottom-0 left-0 right-0 ${cstyles} transition-all ease-in-out`}
      aria-labelledby="modal-title"
      role="dialog"
      aria-modal="true"
    >
      <div className="flex min-h-full items-center justify-center p-4 text-center sm:p-0">
        <div className="relative transform overflow-hidden rounded-lg bg-white text-left shadow-xl transition-all sm:my-8 sm:w-full sm:max-w-lg">
          <div
            className={`${mstyle} bg-white px-4 pt-5 pb-4 sm:p-6 sm:pb-4 flex`}
          >
            <div className="w-1/12 flex justify-start items-start mr-2">
              <svg
                className="h-6 w-6 text-red-600"
                xmlns="http://www.w3.org/2000/svg"
                fill="none"
                viewBox="0 0 24 24"
                strokeWidth="1.5"
                stroke="currentColor"
                aria-hidden="true"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  d="M12 10.5v3.75m-9.303 3.376C1.83 19.126 2.914 21 4.645 21h14.71c1.73 0 2.813-1.874 1.948-3.374L13.949 4.88c-.866-1.501-3.032-1.501-3.898 0L2.697 17.626zM12 17.25h.007v.008H12v-.008z"
                />
              </svg>
            </div>
            <div className="w-10/12 flex flex-col justify-start">
              <h3
                className="text-lg max-h-20 font-medium leading-6 text-gray-900 max-w-max flex justify-center items-center overflow-scroll no-scrollbar"
                id="modal-title"
              >
                {title}
              </h3>
              <div className="mt-2 max-h-40 overflow-scroll no-scrollbar">
                <p className="text-sm text-gray-500">{text}</p>
              </div>
              <div className="flex md:flex-row flex-col justify-between items-center">
                <div className="md:w-1/2 md:pr-2 w-full py-2">
                  <OnlyofficeButton
                    text={mainBtn}
                    onClick={mainAction}
                    primary
                    fullWidth
                  />
                </div>
                <div className="md:w-1/2 md:pl-2 w-full py-2">
                  <OnlyofficeButton
                    text={secBtn}
                    onClick={secAction}
                    fullWidth
                  />
                </div>
              </div>
            </div>
            <div className="w-1/12 flex justify-end items-start">
              <button type="button" onClick={close}>
                <svg
                  width="16"
                  height="16"
                  viewBox="0 0 16 16"
                  fill="none"
                  xmlns="http://www.w3.org/2000/svg"
                >
                  <path
                    fillRule="evenodd"
                    clipRule="evenodd"
                    d="M8.00001 6.58579L3.05027 1.63604L1.63605 3.05025L6.5858 8L1.63605 12.9497L3.05027 14.364L8.00001 9.41421L12.9498 14.364L14.364 12.9497L9.41423 8L14.364 3.05025L12.9498 1.63604L8.00001 6.58579Z"
                    fill="#979797"
                  />
                </svg>
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};
