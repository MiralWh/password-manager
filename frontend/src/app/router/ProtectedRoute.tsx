import { ReactNode } from 'react';
import { Navigate } from 'react-router-dom';
import { authStorage } from '@/shared/lib/storage';

export function ProtectedRoute({ children }: { children: ReactNode }) {
  if (!authStorage.getToken()) return <Navigate to="/login" replace />;
  return <>{children}</>;
}
