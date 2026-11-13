import { Avatar, Box, Card, CardContent, Divider, Link, Stack, Typography } from '@mui/material';
import { Code, Email, GitHub } from '@mui/icons-material';
import { Header } from '@/widgets/header/Header';

export function AboutPage() {
  return (
    <Box sx={{ minHeight: '100vh' }}>
      <Header />
      <Box sx={{ maxWidth: 800, mx: 'auto', p: 3 }}>
        <Typography variant="h5" gutterBottom fontWeight={600}>О создателе приложения</Typography>
        <Card>
          <CardContent>
            <Stack direction="row" spacing={2} alignItems="center">
              <Avatar sx={{ width: 72, height: 72, bgcolor: 'primary.main' }}><Code fontSize="large" /></Avatar>
              <Box>
                <Typography variant="h6">Ваше Имя</Typography>
                <Typography variant="body2" color="text.secondary">Senior Go & React Developer</Typography>
              </Box>
            </Stack>
            <Divider sx={{ my: 2 }} />
            <Typography variant="body1" paragraph>
              Привет! Это демонстрационное приложение Password Manager, разработанное как
              пример клиент-серверной архитектуры с использованием Go (backend) и
              React + Material UI (frontend).
            </Typography>
            <Typography variant="body1" paragraph>
              Backend реализован на Go с чистой архитектурой, использует PostgreSQL
              для хранения данных, AES-256-GCM для шифрования паролей и JWT для аутентификации.
            </Typography>
            <Typography variant="body1" paragraph>
              Frontend построен по методологии Feature-Sliced Design (FSD) с использованием
              React, TypeScript и библиотеки Material UI.
            </Typography>
            <Divider sx={{ my: 2 }} />
            <Stack spacing={1}>
              <Stack direction="row" spacing={1} alignItems="center">
                <GitHub fontSize="small" />
                <Link href="https://github.com/" target="_blank" rel="noopener">github.com/your-username</Link>
              </Stack>
              <Stack direction="row" spacing={1} alignItems="center">
                <Email fontSize="small" />
                <Link href="mailto:you@example.com">you@example.com</Link>
              </Stack>
            </Stack>
          </CardContent>
        </Card>
      </Box>
    </Box>
  );
}
