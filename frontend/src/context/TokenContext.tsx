/* eslint-disable @typescript-eslint/no-empty-function */
import React, { useContext, useRef, useCallback } from "react";
import axios from "axios";
import zoomSdk from "@zoom/appssdk";
import { UserResponse } from "src/types/user";

type ProviderProps = {
  children?: JSX.Element | JSX.Element[];
};

type NullableUserResponse = UserResponse | null;

const TokenContext = React.createContext<() => Promise<NullableUserResponse>>(
  () =>
    new Promise((res) => {
      res(null);
    })
);

export const TokenProvider: React.FC<ProviderProps> = ({ children }) => {
  const token = useRef<UserResponse | null>(null);
  const updateToken = useCallback(async (): Promise<NullableUserResponse> => {
    if (token.current && Date.now() < token.current.expires_at) {
      return token.current;
    }

    try {
      const zctx = await zoomSdk.getAppContext();
      const resp = await axios.get<UserResponse>(
        `${process.env.BACKEND_GATEWAY}/api/me`,
        {
          timeout: 2500,
          headers: {
            "X-Zoom-App-Context": zctx.context,
          },
        }
      );
      token.current = resp.data;
      return token.current;
    } catch {
      return null;
    }
  }, []);

  return (
    <TokenContext.Provider value={updateToken}>
      {children}
    </TokenContext.Provider>
  );
};

export const useToken = () => useContext(TokenContext);
