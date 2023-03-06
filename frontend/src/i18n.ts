import i18n from "i18next";
import I18NextHttpBackend from "i18next-http-backend";
import LanguageDetector from "i18next-browser-languagedetector";
import { initReactI18next } from "react-i18next";

i18n
  .use(LanguageDetector)
  .use(I18NextHttpBackend)
  .use(initReactI18next)
  .init({
    fallbackLng: "en-US",
    debug: false,
    interpolation: {
      escapeValue: false,
    },
    detection: {
      order: ["navigator"],
    },
  });

export default i18n;
