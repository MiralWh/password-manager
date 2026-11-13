import { FormEvent, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  Alert, Box, Button, CircularProgress, IconButton, InputAdornment, TextField, Typography,
} from '@mui/material';
import { Visibility, VisibilityOff } from '@mui/icons-material';
import { authApi } from '@/entities/user/api';
import { authStorage } from '@/shared/lib/storage';
import { getErrorMessage } from '@/shared/api/client';

export function LoginForm() {
  const navigate = useNavigate();
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [show, setShow] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  const submit = async (e: FormEvent) => {
    e.preventDefault();
    setError(null);
    setLoading(true);
    try {
      const res = await authApi.login(username, password);
      authStorage.setSession(res.token, res.username, res.expires_at);
      navigate('/passwords');
    } catch (err) {
      setError(getErrorMessage(err));
    } finally {
      setLoading(false);
    }
  };

  return (
    <Box component="form" onSubmit={submit} sx={{ display: 'flex', flexDirection: 'column', gap: 2 }}>
      <Typography variant="h4" align="center" fontWeight={600}>Вход в Password Manager</Typography>
      {error && <Alert severity="error">{error}</Alert>}
      <TextField label="Логин" value={username} onChange={(e) => setUsername(e.target.value)} required autoFocus fullWidth />
      <TextField
        label="Пароль"
        type={show ? 'text' : 'password'}
        value={password}
        onChange={(e) => setPassword(e.target.value)}
        required
        fullWidth
        InputProps={{
          endAdornment: (
            <InputAdornment position="end">
              <IconButton onClick={() => setShow((v) => !v)} edge="end" type="button">
                {show ? <VisibilityOff /> : <Visibility />}
              </IconButton>
            </InputAdornment>
          ),
        }}
      />
      <Button type="submit" variant="contained" size="large" disabled={loading}
        startIcon={loading ? <CircularProgress size={20} color="inherit" /> : null}>
        Войти
      </Button>
      <Typography variant="body2" color="text.secondary" align="center">
        По умолчанию: admin / admin123
      </Typography>
    </Box>
  );
}
