import { Navigate, Route, Routes } from 'react-router-dom';
import { LoginPage } from '@/pages/login/LoginPage';
import { PasswordsPage } from '@/pages/passwords/PasswordsPage';
import { AboutPage } from '@/pages/about/AboutPage';
import { SettingsPage } from '@/pages/settings/SettingsPage';
import { ProtectedRoute } from './ProtectedRoute';

export function AppRouter() {
  return (
    <Routes>
      <Route path="/login" element={<LoginPage />} />
      <Route path="/passwords" element={<ProtectedRoute><PasswordsPage /></ProtectedRoute>} />
      <Route path="/settings" element={<ProtectedRoute><SettingsPage /></ProtectedRoute>} />
      <Route path="/about" element={<ProtectedRoute><AboutPage /></ProtectedRoute>} />
      <Route path="/" element={<Navigate to="/passwords" replace />} />
      <Route path="*" element={<Navigate to="/passwords" replace />} />
    </Routes>
  );
}
