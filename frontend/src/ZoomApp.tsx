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
  useNavigate,
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

import icon from "@assets/nofile.svg";

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
  const navigate = useNavigate();
  const { ready, error } = useMainContext();
  const { error: socketError } = useSnapshot(SocketState);
  if (ready && !location.pathname.includes("editor") && (error || socketError))
    return (
      <OnlyofficeBasicLayoutContainer>
        <React.Suspense fallback={<CenteredOnlyofficeSpinner />}>
          <div className="w-screen h-screen flex justify-center flex-col items-center mb-1">
            <img src={icon} alt="error-icon" />
            <OnlyofficeError
              text={
                t("context.error") ||
                "Could not fetch user information or establish a new websocket connection"
              }
            />
            <div className="pt-5">
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
              <OnlyofficeBasicLayoutContainer
                onNavbarClick={() => navigate("/")}
              >
                <React.Suspense fallback={<CenteredOnlyofficeSpinner />}>
                  <MainPage />
                </React.Suspense>
              </OnlyofficeBasicLayoutContainer>
            }
          />
          <Route
            path="create"
            element={
              <OnlyofficeBasicLayoutContainer
                onNavbarClick={() => navigate("/")}
              >
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
