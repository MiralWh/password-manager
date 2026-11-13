import { AppBar, Box, Button, Toolbar, Typography } from '@mui/material';
import { Info, Lock, Logout, Settings, VpnKey } from '@mui/icons-material';
import { useLocation, useNavigate } from 'react-router-dom';
import { authStorage } from '@/shared/lib/storage';

export function Header() {
  const navigate = useNavigate();
  const { pathname } = useLocation();
  const username = authStorage.getUsername();

  const logout = () => {
    authStorage.clear();
    navigate('/login', { replace: true });
  };

  const btnSx = (active: boolean) => ({
    backgroundColor: active ? 'rgba(255,255,255,0.15)' : 'transparent',
  });

  return (
    <AppBar position="static" color="primary" elevation={1}>
      <Toolbar>
        <Lock sx={{ mr: 1 }} />
        <Typography variant="h6" sx={{ flexGrow: 1 }}>Password Manager</Typography>
        <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
          <Button color="inherit" startIcon={<VpnKey />} onClick={() => navigate('/passwords')}
            sx={btnSx(pathname === '/passwords')}>Пароли</Button>
          <Button color="inherit" startIcon={<Settings />} onClick={() => navigate('/settings')}
            sx={btnSx(pathname === '/settings')}>Настройки</Button>
          <Button color="inherit" startIcon={<Info />} onClick={() => navigate('/about')}
            sx={btnSx(pathname === '/about')}>О создателе</Button>
          {username && <Typography variant="body2" sx={{ mx: 1, opacity: 0.9 }}>{username}</Typography>}
          <Button color="inherit" startIcon={<Logout />} onClick={logout}>Выйти</Button>
        </Box>
      </Toolbar>
    </AppBar>
  );
}
