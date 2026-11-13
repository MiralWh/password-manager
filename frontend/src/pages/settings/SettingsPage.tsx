import { useEffect, useState } from 'react';
import {
  Alert, Box, Button, Card, CardContent, CircularProgress, Dialog, DialogActions,
  DialogContent, DialogContentText, DialogTitle, Divider, IconButton, Paper,
  Snackbar, Stack, Table, TableBody, TableCell, TableHead, TableRow, TextField, Tooltip, Typography,
} from '@mui/material';
import { Add, Delete, Key, OpenInNew } from '@mui/icons-material';
import { Header } from '@/widgets/header/Header';
import { apiKeyApi } from '@/entities/apikey/api';
import { APIKey, APIKeyCreated } from '@/entities/apikey/model';
import { CopyButton } from '@/shared/ui/CopyButton';
import { getErrorMessage } from '@/shared/api/client';

export function SettingsPage() {
  const [keys, setKeys] = useState<APIKey[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [snack, setSnack] = useState<string | null>(null);
  const [createOpen, setCreateOpen] = useState(false);
  const [newName, setNewName] = useState('');
  const [creating, setCreating] = useState(false);
  const [createdKey, setCreatedKey] = useState<APIKeyCreated | null>(null);
  const [deleteId, setDeleteId] = useState<number | null>(null);

  useEffect(() => { void load(); }, []);
  const load = async () => {
    setLoading(true);
    try { setKeys(await apiKeyApi.list()); }
    catch (err) { setError(getErrorMessage(err)); }
    finally { setLoading(false); }
  };

  const create = async () => {
    if (!newName.trim()) return;
    setCreating(true);
    try {
      const k = await apiKeyApi.create(newName.trim());
      setCreatedKey(k);
      setCreateOpen(false);
      setNewName('');
      await load();
    } catch (err) { setSnack(getErrorMessage(err)); }
    finally { setCreating(false); }
  };

  const remove = async () => {
    if (deleteId == null) return;
    try {
      await apiKeyApi.remove(deleteId);
      setSnack('Ключ отозван');
      await load();
    } catch (err) { setSnack(getErrorMessage(err)); }
    finally { setDeleteId(null); }
  };

  return (
    <Box sx={{ minHeight: '100vh' }}>
      <Header />
      <Box sx={{ maxWidth: 960, mx: 'auto', p: 3 }}>
        <Typography variant="h5" fontWeight={600} gutterBottom>Настройки</Typography>

        <Card sx={{ mb: 3 }}>
          <CardContent>
            <Stack direction="row" alignItems="center" spacing={1} sx={{ mb: 2 }}>
              <OpenInNew color="primary" />
              <Typography variant="h6" sx={{ flexGrow: 1 }}>API для десктопа</Typography>
              <Button
                variant="outlined" size="small"
                href="/swagger/" target="_blank" rel="noopener"
              >
                Swagger UI
              </Button>
            </Stack>
            <Typography variant="body2" color="text.secondary" paragraph>
              Используйте <b>X-API-Key</b> в заголовке запроса вместо JWT. Полезно для десктопных
              приложений и скриптов. Base URL: <code>/api/v1</code>. Ключ показывается один раз
              при создании — сохраните его.
            </Typography>

            <Divider sx={{ my: 2 }} />

            <Stack direction="row" alignItems="center" spacing={1} sx={{ mb: 2 }}>
              <Key color="primary" />
              <Typography variant="subtitle1" sx={{ flexGrow: 1 }}>API-ключи</Typography>
              <Button variant="contained" startIcon={<Add />} onClick={() => setCreateOpen(true)}>
                Создать
              </Button>
            </Stack>

            {error && <Alert severity="error" sx={{ mb: 2 }}>{error}</Alert>}

            {loading ? (
              <Box sx={{ display: 'flex', justifyContent: 'center', py: 4 }}><CircularProgress /></Box>
            ) : keys.length === 0 ? (
              <Paper sx={{ p: 3, textAlign: 'center' }} variant="outlined">
                <Typography color="text.secondary">Пока нет API-ключей.</Typography>
              </Paper>
            ) : (
              <Table size="small">
                <TableHead>
                  <TableRow>
                    <TableCell>Название</TableCell>
                    <TableCell>Префикс</TableCell>
                    <TableCell>Создан</TableCell>
                    <TableCell>Последнее использование</TableCell>
                    <TableCell align="right"></TableCell>
                  </TableRow>
                </TableHead>
                <TableBody>
                  {keys.map((k) => (
                    <TableRow key={k.id}>
                      <TableCell>{k.name}</TableCell>
                      <TableCell><code>{k.key_prefix}…</code></TableCell>
                      <TableCell>{new Date(k.created_at).toLocaleString()}</TableCell>
                      <TableCell>
                        {k.last_used_at ? new Date(k.last_used_at).toLocaleString() : '—'}
                      </TableCell>
                      <TableCell align="right">
                        <Tooltip title="Отозвать">
                          <IconButton size="small" color="error" onClick={() => setDeleteId(k.id)}>
                            <Delete fontSize="small" />
                          </IconButton>
                        </Tooltip>
                      </TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            )}
          </CardContent>
        </Card>
      </Box>

      <Dialog open={createOpen} onClose={() => setCreateOpen(false)} fullWidth maxWidth="xs">
        <DialogTitle>Новый API-ключ</DialogTitle>
        <DialogContent>
          <TextField
            autoFocus fullWidth margin="dense" label="Название (например, Desktop app)"
            value={newName} onChange={(e) => setNewName(e.target.value)}
            onKeyDown={(e) => { if (e.key === 'Enter') void create(); }}
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setCreateOpen(false)}>Отмена</Button>
          <Button variant="contained" disabled={creating} onClick={create}>Создать</Button>
        </DialogActions>
      </Dialog>

      <Dialog open={!!createdKey} onClose={() => setCreatedKey(null)} fullWidth maxWidth="sm">
        <DialogTitle>API-ключ создан</DialogTitle>
        <DialogContent>
          <Alert severity="warning" sx={{ mb: 2 }}>
            Это <b>единственный раз</b>, когда вы видите полный ключ. Скопируйте и сохраните.
          </Alert>
          <Paper variant="outlined" sx={{ p: 2, display: 'flex', alignItems: 'center', gap: 1 }}>
            <Typography sx={{ fontFamily: 'monospace', wordBreak: 'break-all', flexGrow: 1 }}>
              {createdKey?.key}
            </Typography>
            <CopyButton value={createdKey?.key ?? ''} size="medium" tooltip="Копировать ключ" />
          </Paper>
        </DialogContent>
        <DialogActions>
          <Button variant="contained" onClick={() => setCreatedKey(null)}>Я сохранил ключ</Button>
        </DialogActions>
      </Dialog>

      <Dialog open={deleteId != null} onClose={() => setDeleteId(null)}>
        <DialogTitle>Отозвать ключ?</DialogTitle>
        <DialogContent>
          <DialogContentText>Приложения, использующие этот ключ, потеряют доступ.</DialogContentText>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setDeleteId(null)}>Отмена</Button>
          <Button color="error" variant="contained" onClick={remove}>Отозвать</Button>
        </DialogActions>
      </Dialog>

      <Snackbar open={!!snack} autoHideDuration={3000} onClose={() => setSnack(null)} message={snack ?? ''} />
    </Box>
  );
}
