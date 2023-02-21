import React from "react";
import { t } from "i18next";
import { QueryClient, QueryClientProvider } from "react-query";
import { useSnapshot } from "valtio";
import {
  BrowserRouter as Router,
  Routes,
  Route,
  Navigate,
  useLocation,
} from "react-router-dom";

import { OnlyofficeBasicLayoutContainer } from "@layouts/container";

import { OnlyofficeSpinner } from "@components/spinner";

import {
  MainProvider,
  SocketState,
  useMainContext,
} from "@context/MainContext";
import { OnlyofficeError } from "@components/error";
import { OnlyofficeButton } from "@components/button";

import BackgroundError from "@assets/background-error.svg";
import { OnlyofficeSubtitle } from "@components/title";

const MainPage = React.lazy(() => import("@pages/Main"));
const CreationPage = React.lazy(() => import("@pages/Creation"));
const OnlyofficeEditorPage = React.lazy(() => import("@pages/Editor"));
const CenteredOnlyofficeSpinner = () => (
  <div className="w-full h-full flex justify-center items-center">
    <OnlyofficeSpinner />
  </div>
);

const queryClient = new QueryClient();

const LazyRoutes: React.FC = () => {
  const location = useLocation();
  const { ready, error } = useMainContext();
  const { error: socketError } = useSnapshot(SocketState);
  if (ready && !location.pathname.includes("editor") && (error || socketError))
    return (
      <OnlyofficeBasicLayoutContainer>
        <React.Suspense fallback={<CenteredOnlyofficeSpinner />}>
          <div className="w-screen h-screen flex justify-center flex-col items-center mb-1">
            <div className="absolute flex justify-center items-center w-screen h-screen">
              <BackgroundError />
            </div>
            <div className="pb-5">
              <OnlyofficeError text={t("context.error.title") || "Error"} />
            </div>
            <OnlyofficeSubtitle
              text={
                t("context.error.text") ||
                "Something went wrong. Please reload the page or contact the site administrator."
              }
            />
            <div className="pt-5 z-[100]">
              <OnlyofficeButton
                primary
                text={t("button.reload") || "Reload"}
                onClick={() => window.location.reload()}
              />
            </div>
          </div>
        </React.Suspense>
      </OnlyofficeBasicLayoutContainer>
    );

  if (ready && !error && (!socketError || location.pathname.includes("editor")))
    return (
      <Routes location={location} key={location.pathname}>
        <Route path="/">
          <Route
            index
            element={
              <OnlyofficeBasicLayoutContainer>
                <React.Suspense fallback={<CenteredOnlyofficeSpinner />}>
                  <MainPage />
                </React.Suspense>
              </OnlyofficeBasicLayoutContainer>
            }
          />
          <Route
            path="create"
            element={
              <OnlyofficeBasicLayoutContainer>
                <React.Suspense fallback={<CenteredOnlyofficeSpinner />}>
                  <CreationPage />
                </React.Suspense>
              </OnlyofficeBasicLayoutContainer>
            }
          />
          <Route
            path="editor"
            element={
              <React.Suspense fallback={<OnlyofficeSpinner />}>
                <OnlyofficeEditorPage />
              </React.Suspense>
            }
          />
        </Route>
        <Route path="*" element={<Navigate to="/" />} />
      </Routes>
    );

  return (
    <OnlyofficeBasicLayoutContainer>
      <CenteredOnlyofficeSpinner />
    </OnlyofficeBasicLayoutContainer>
  );
};

const ZoomApp: React.FC = () => (
  <MainProvider>
    <QueryClientProvider client={queryClient}>
      <Router>
        <LazyRoutes />
      </Router>
    </QueryClientProvider>
  </MainProvider>
);

export default ZoomApp;
