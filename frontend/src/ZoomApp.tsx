import React from "react";
import { QueryClient, QueryClientProvider } from "react-query";
import {
  BrowserRouter as Router,
  Routes,
  Route,
  Navigate,
  useLocation,
  useNavigate,
} from "react-router-dom";

import { OnlyofficeSpinner } from "@components/spinner";

import { OnlyofficeBasicLayoutContainer } from "@layouts/container";

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
  return (
    <Routes location={location} key={location.pathname}>
      <Route path="/">
        <Route
          index
          element={
            <OnlyofficeBasicLayoutContainer onNavbarClick={() => navigate("/")}>
              <React.Suspense fallback={<CenteredOnlyofficeSpinner />}>
                <MainPage />
              </React.Suspense>
            </OnlyofficeBasicLayoutContainer>
          }
        />
        <Route
          path="create"
          element={
            <OnlyofficeBasicLayoutContainer onNavbarClick={() => navigate("/")}>
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
};

function ZoomApp() {
  return (
    <QueryClientProvider client={queryClient}>
      <Router>
        <LazyRoutes />
      </Router>
    </QueryClientProvider>
  );
}

export default ZoomApp;
