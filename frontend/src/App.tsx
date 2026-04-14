import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import { AppLayout } from './components/layout/AppLayout'
import { LoginPage } from './pages/LoginPage'
import { DashboardPage } from './pages/DashboardPage'
import { TicketListPage } from './pages/TicketListPage'
import { TicketDetailPage } from './pages/TicketDetailPage'
import { TicketFormPage } from './pages/TicketFormPage'
import { ApprovalsPage } from './pages/ApprovalsPage'
import { ActivityLogsPage } from './pages/ActivityLogsPage'
import { NotificationsPage } from './pages/NotificationsPage'
import { UserManagementPage } from './pages/UserManagementPage'
import { WebhookConfigPage } from './pages/WebhookConfigPage'
import { OrgStructurePage } from './pages/OrgStructurePage'

import { SSOCallbackPage } from './pages/SSOCallbackPage'

export default function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/login" element={<LoginPage />} />
        <Route path="/sso/callback" element={<SSOCallbackPage />} />
        <Route element={<AppLayout />}>
          <Route path="/" element={<Navigate to="/dashboard" replace />} />
          <Route path="/dashboard" element={<DashboardPage />} />
          <Route path="/tickets" element={<TicketListPage />} />
          <Route path="/tickets/new" element={<TicketFormPage />} />
          <Route path="/tickets/:id" element={<TicketDetailPage />} />
          <Route path="/approvals" element={<ApprovalsPage />} />
          <Route path="/activity-logs" element={<ActivityLogsPage />} />
          <Route path="/notifications" element={<NotificationsPage />} />
          <Route path="/users" element={<UserManagementPage />} />
          <Route path="/webhooks" element={<WebhookConfigPage />} />
          <Route path="/org-structure" element={<OrgStructurePage />} />
          <Route path="/settings" element={<UserManagementPage />} />
        </Route>
      </Routes>
    </BrowserRouter>
  )
}
