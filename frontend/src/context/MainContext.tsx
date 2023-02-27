import React, {
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useRef,
  useState,
} from "react";
import ReconnectingWebSocket from "reconnecting-websocket";
import { proxy } from "valtio";
import md5 from "md5";
import zoomSdk from "@zoom/appssdk";

import { getMe } from "@services/me";

import { MeResponse } from "src/types/user";

export const CurrentUser = proxy<MeResponse>({
  id: "",
  first_name: "",
  last_name: "",
  email: "",
  language: "en-US",
});

export const SocketState = proxy<{
  ready: boolean;
  error: boolean;
  value: any;
}>({
  ready: false,
  error: false,
  value: null,
});

type ProviderProps = {
  children?: JSX.Element | JSX.Element[];
};

const MainContext = React.createContext<{
  ready: boolean;
  error: boolean;
}>({ ready: false, error: false });

export const MainProvider: React.FC<ProviderProps> = ({ children }) => {
  const [ready, setReady] = useState(false);
  const [error, setError] = useState(false);
  const rws = useRef<ReconnectingWebSocket | null>(null);
  const res = useMemo(
    () => ({
      ready,
      error,
    }),
    [ready, error]
  );

  const urlProvider = useCallback(async () => {
    try {
      const meeting = await zoomSdk.getMeetingUUID();
      const appCtx = await zoomSdk.getAppContext();
      return `${new URL(process.env.BACKEND_GATEWAY_WS || "").toString()}/${md5(
        meeting.meetingUUID
      )}?token=${appCtx.context}`;
    } catch {
      return "";
    }
  }, []);

  useEffect(() => {
    let socket: ReconnectingWebSocket | null;
    setReady(false);
    setError(false);
    Promise.all([
      getMe().then((result) => {
        const {
          id,
          email,
          first_name: firstName,
          last_name: lastName,
          language,
        } = result.response;
        CurrentUser.id = id;
        CurrentUser.email = email;
        CurrentUser.first_name = firstName;
        CurrentUser.last_name = lastName;
        CurrentUser.language = language;
      }),
      zoomSdk
        .getMeetingUUID()
        .then(() => {
          socket = new ReconnectingWebSocket(urlProvider, ["wss", "ws"], {
            maxRetries: 5,
          });
          socket.onopen = () => {
            SocketState.ready = true;
            SocketState.error = false;
          };
          socket.onclose = () => {
            SocketState.ready = false;
          };
          socket.onmessage = (event) => {
            SocketState.value = event.data;
          };
          socket.onerror = () => {
            if (socket?.retryCount === 5) {
              SocketState.error = true;
            }
          };
          rws.current = socket;
        })
        .catch(() => {
          SocketState.ready = true;
        }),
    ])
      .then(() => {
        setReady(true);
        setError(false);
      })
      .catch(() => {
        setReady(true);
        setError(true);
      });
    return () => {
      socket?.close();
    };
  }, [urlProvider]);
  return <MainContext.Provider value={res}>{children}</MainContext.Provider>;
};

export const useMainContext = () => useContext(MainContext);
