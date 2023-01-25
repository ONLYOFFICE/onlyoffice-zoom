/* eslint-disable react/jsx-no-constructed-context-values */
import React, {
  useState,
  useEffect,
  useRef,
  createContext,
  useContext,
  useCallback,
} from "react";
import md5 from "md5";
import zoomSdk from "@zoom/appssdk";
import ReconnectingWebSocket from "reconnecting-websocket";

const WebsocketContext = createContext<{
  ready: boolean;
  error: boolean;
  value?: any;
}>({ ready: false, error: false });

export const WebsocketProvider: React.FC<{
  children: React.ReactNode;
}> = ({ children }) => {
  const [ready, setReady] = useState(false);
  const [error, setError] = useState(false);
  const [value, setValue] = useState(null);
  const rws = useRef<ReconnectingWebSocket | null>(null);
  const urlProvider = useCallback(async () => {
    const meeting = await zoomSdk.getMeetingUUID();
    const appCtx = await zoomSdk.getAppContext();
    return `${new URL(process.env.BACKEND_GATEWAY_WS || "").toString()}/${md5(
      meeting.meetingUUID
    )}?token=${appCtx.context}`;
  }, []);

  useEffect(() => {
    const socket = new ReconnectingWebSocket(urlProvider);
    socket.onopen = () => {
      setReady(true);
      setError(false);
    };
    socket.onclose = () => setReady(false);
    socket.onmessage = (event) => setValue(event.data);
    socket.onerror = () => setError(true);
    rws.current = socket;
    return () => {
      socket.close();
    };
  }, [urlProvider]);

  return (
    <WebsocketContext.Provider
      value={{
        ready,
        error,
        value,
      }}
    >
      {children}
    </WebsocketContext.Provider>
  );
};

export const useWebsocket = () => useContext(WebsocketContext);
