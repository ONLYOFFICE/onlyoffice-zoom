import React, {
  useState,
  useEffect,
  createContext,
  useContext,
  useMemo,
} from "react";
import { getMe } from "@services/me";

const LangContext = createContext<{ lang: string; loading: boolean }>({
  lang: "en-US",
  loading: true,
});

export const LangProvider: React.FC<{
  children: React.ReactNode;
}> = ({ children }) => {
  const [lang, setLang] = useState("en-US");
  const [loading, setLoading] = useState(true);
  const value = useMemo(
    () => ({
      lang,
      loading,
    }),
    [lang, loading]
  );

  useEffect(() => {
    setLoading(true);
    getMe()
      .then((res) => {
        setLang(res.response.language || "en-US");
      })
      .catch(() => setLang("en-US"))
      .finally(() => setLoading(false));
  }, []);

  return <LangContext.Provider value={value}>{children}</LangContext.Provider>;
};

export const useZoomLanguage = () => useContext(LangContext);
