import React, {
  useState,
  useEffect,
  createContext,
  useContext,
  useMemo,
} from "react";
import i18n from "i18next";
import { getMe } from "@services/me";

const LangContext = createContext<{ loading: boolean }>({
  loading: true,
});

export const LangProvider: React.FC<{
  children: React.ReactNode;
}> = ({ children }) => {
  const [loading, setLoading] = useState(true);
  const value = useMemo(
    () => ({
      loading,
    }),
    [loading]
  );

  useEffect(() => {
    setLoading(true);
    getMe()
      .then((res) => i18n.changeLanguage(res.response.language || "en-US"))
      .catch(() => i18n.changeLanguage("en-US"))
      .finally(() => setLoading(false));
  }, []);

  return <LangContext.Provider value={value}>{children}</LangContext.Provider>;
};

export const useZoomLanguage = () => useContext(LangContext);
