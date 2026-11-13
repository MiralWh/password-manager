import { useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { Box, Paper } from '@mui/material';
import { LoginForm } from '@/features/auth/login/LoginForm';
import { authStorage } from '@/shared/lib/storage';

export function LoginPage() {
  const navigate = useNavigate();
  useEffect(() => {
    if (authStorage.getToken()) navigate('/passwords', { replace: true });
  }, [navigate]);

  return (
    <Box sx={{
      minHeight: '100vh', display: 'flex', alignItems: 'center', justifyContent: 'center',
      background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)', p: 2,
    }}>
      <Paper elevation={6} sx={{ p: 4, maxWidth: 420, width: '100%', borderRadius: 3 }}>
        <LoginForm />
      </Paper>
    </Box>
  );
}
