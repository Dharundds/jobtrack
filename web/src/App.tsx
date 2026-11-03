import { BrowserRouter, Navigate, Route, Routes } from "react-router-dom";
import { AppLayout } from "@/components/layout/AppLayout";
import { RequireConnected } from "@/components/RequireConnected";
import { ThemeProvider } from "@/lib/theme";
import { ApplicationDetailPage } from "@/pages/ApplicationDetailPage";
import { ApplicationsListPage } from "@/pages/ApplicationsListPage";
import { BoardPage } from "@/pages/BoardPage";
import { ConnectPage } from "@/pages/ConnectPage";
import { DashboardPage } from "@/pages/DashboardPage";
import { SettingsPage } from "@/pages/SettingsPage";

export default function App() {
  return (
    <ThemeProvider>
      <BrowserRouter>
        <Routes>
          <Route path="/connect" element={<ConnectPage />} />
          <Route element={<RequireConnected />}>
            <Route element={<AppLayout />}>
              <Route index element={<DashboardPage />} />
              <Route path="board" element={<BoardPage />} />
              <Route path="applications" element={<ApplicationsListPage />} />
              <Route path="applications/:id" element={<ApplicationDetailPage />} />
              <Route path="settings" element={<SettingsPage />} />
            </Route>
          </Route>
          <Route path="*" element={<Navigate to="/board" replace />} />
        </Routes>
      </BrowserRouter>
    </ThemeProvider>
  );
}
