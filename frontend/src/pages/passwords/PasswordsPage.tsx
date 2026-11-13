import { useEffect, useMemo, useState } from 'react';
import {
  Alert, Box, Button, Card, CardActions, CardContent, Chip, CircularProgress,
  Dialog, DialogActions, DialogContent, DialogContentText, DialogTitle,
  IconButton, InputAdornment, Paper, Snackbar, TextField, Tooltip, Typography,
} from '@mui/material';
import { Add, Clear, Delete, Edit, Search } from '@mui/icons-material';
import { Header } from '@/widgets/header/Header';
import { PasswordForm } from '@/widgets/password-form/PasswordForm';
import { GroupsBar } from '@/widgets/groups-bar/GroupsBar';
import { ToggleVisibilityButton } from '@/features/passwords/toggle-visibility/ToggleVisibilityButton';
import { CopyButton } from '@/shared/ui/CopyButton';
import { passwordApi } from '@/entities/password/api';
import { groupApi } from '@/entities/group/api';
import { PasswordEntry, PasswordInput } from '@/entities/password/model';
import { Group } from '@/entities/group/model';
import { getErrorMessage } from '@/shared/api/client';

export function PasswordsPage() {
  const [items, setItems] = useState<PasswordEntry[]>([]);
  const [groups, setGroups] = useState<Group[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [visibleIds, setVisibleIds] = useState<Set<number>>(new Set());
  const [selectedGroup, setSelectedGroup] = useState<number | null>(null);

  // search
  const [search, setSearch] = useState('');
  const [debouncedSearch, setDebouncedSearch] = useState('');

  // dialogs
  const [dialogOpen, setDialogOpen] = useState(false);
  const [editing, setEditing] = useState<PasswordEntry | null>(null);
  const [submitting, setSubmitting] = useState(false);
  const [deleteId, setDeleteId] = useState<number | null>(null);
  const [snack, setSnack] = useState<string | null>(null);
  const [groupDialog, setGroupDialog] = useState<{ open: boolean; name: string }>({ open: false, name: '' });
  const [deleteGroupId, setDeleteGroupId] = useState<number | null>(null);

  // debounce
  useEffect(() => {
    const t = setTimeout(() => setDebouncedSearch(search.trim()), 300);
    return () => clearTimeout(t);
  }, [search]);

  // reload on search change
  useEffect(() => {
    void load();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [debouncedSearch]);

  // при активном поиске сбрасываем выбор группы
  useEffect(() => {
    if (debouncedSearch) setSelectedGroup(null);
  }, [debouncedSearch]);

  const load = async () => {
    setLoading(true);
    setError(null);
    try {
      const [ps, gs] = await Promise.all([
        passwordApi.list(debouncedSearch || undefined),
        groupApi.list(),
      ]);
      setItems(ps);
      setGroups(gs);
    } catch (err) { setError(getErrorMessage(err)); }
    finally { setLoading(false); }
  };

  // группы для отображения — фильтруем клиентом по строке поиска
  const visibleGroups = useMemo(() => {
    if (!debouncedSearch) return groups;
    const s = debouncedSearch.toLowerCase();
    return groups.filter((g) => g.name.toLowerCase().includes(s));
  }, [groups, debouncedSearch]);

  const filtered = useMemo(
    () => selectedGroup === null ? items : items.filter((i) => i.group_id === selectedGroup),
    [items, selectedGroup],
  );

  const toggleVisibility = (id: number) => {
    setVisibleIds((prev) => {
      const next = new Set(prev);
      next.has(id) ? next.delete(id) : next.add(id);
      return next;
    });
  };

  const refreshGroups = async () => setGroups(await groupApi.list());

  const handleCreateOrUpdate = async (input: PasswordInput) => {
    setSubmitting(true);
    try {
      if (editing) {
        const updated = await passwordApi.update(editing.id, input);
        setItems((prev) => prev.map((it) => (it.id === updated.id ? updated : it)));
        setSnack('Пароль обновлён');
      } else {
        const created = await passwordApi.create(input);
        setItems((prev) => [created, ...prev]);
        setSnack('Пароль создан');
      }
      setDialogOpen(false);
      setEditing(null);
      await refreshGroups();
    } catch (err) { setSnack(getErrorMessage(err)); }
    finally { setSubmitting(false); }
  };

  const handleDelete = async () => {
    if (deleteId == null) return;
    try {
      await passwordApi.remove(deleteId);
      setItems((prev) => prev.filter((it) => it.id !== deleteId));
      setSnack('Пароль удалён');
      await refreshGroups();
    } catch (err) { setSnack(getErrorMessage(err)); }
    finally { setDeleteId(null); }
  };

  const createGroup = async () => {
    const name = groupDialog.name.trim();
    if (!name) return;
    try {
      const g = await groupApi.create({ name });
      setSnack(`Группа «${g.name}» создана`);
      setGroupDialog({ open: false, name: '' });
      await refreshGroups();
    } catch (err) { setSnack(getErrorMessage(err)); }
  };

  const removeGroup = async () => {
    if (deleteGroupId == null) return;
    try {
      await groupApi.remove(deleteGroupId);
      if (selectedGroup === deleteGroupId) setSelectedGroup(null);
      setSnack('Группа удалена');
      await Promise.all([refreshGroups(), load()]);
    } catch (err) { setSnack(getErrorMessage(err)); }
    finally { setDeleteGroupId(null); }
  };

  const getGroupName = (id?: number | null) =>
    id == null ? null : groups.find((g) => g.id === id)?.name ?? null;

  const searchActive = debouncedSearch.length > 0;

  return (
    <Box sx={{ minHeight: '100vh' }}>
      <Header />
      <Box sx={{ maxWidth: 960, mx: 'auto', p: 3 }}>
        <Box sx={{ display: 'flex', alignItems: 'center', mb: 2 }}>
          <Typography variant="h5" sx={{ flexGrow: 1 }} fontWeight={600}>Мои пароли</Typography>
          <Button variant="contained" startIcon={<Add />} onClick={() => { setEditing(null); setDialogOpen(true); }}>
            Добавить
          </Button>
        </Box>

        {/* Search */}
        <TextField
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          placeholder="Поиск: название, логин, URL, заметки, группа"
          fullWidth
          size="small"
          sx={{ mb: 2 }}
          InputProps={{
            startAdornment: (
              <InputAdornment position="start"><Search fontSize="small" /></InputAdornment>
            ),
            endAdornment: search ? (
              <InputAdornment position="end">
                <IconButton size="small" onClick={() => setSearch('')} aria-label="Очистить">
                  <Clear fontSize="small" />
                </IconButton>
              </InputAdornment>
            ) : null,
          }}
        />

        {error && <Alert severity="error" sx={{ mb: 2 }}>{error}</Alert>}

        <GroupsBar
          groups={visibleGroups}
          selectedId={selectedGroup}
          onSelect={setSelectedGroup}
          onCreate={() => setGroupDialog({ open: true, name: '' })}
          onDelete={(id) => setDeleteGroupId(id)}
        />

        {searchActive && !loading && (
          <Typography variant="body2" color="text.secondary" sx={{ mb: 1 }}>
            Найдено: <b>{filtered.length}</b> {filtered.length === 1 ? 'запись' : 'записей'}
            {selectedGroup !== null && getGroupName(selectedGroup) && (
              <> в группе «{getGroupName(selectedGroup)}»</>
            )}
            {' '}по запросу «{debouncedSearch}»
          </Typography>
        )}

        {loading ? (
          <Box sx={{ display: 'flex', justifyContent: 'center', py: 6 }}><CircularProgress /></Box>
        ) : filtered.length === 0 ? (
          <Paper sx={{ p: 4, textAlign: 'center' }}>
            <Typography color="text.secondary">
              {searchActive
                ? `По запросу «${debouncedSearch}» ничего не найдено.`
                : selectedGroup === null
                  ? 'Пока нет сохранённых паролей. Нажмите «Добавить», чтобы создать первый.'
                  : 'В этой группе пока нет паролей.'}
            </Typography>
          </Paper>
        ) : (
          <Box sx={{ display: 'grid', gridTemplateColumns: { xs: '1fr', sm: '1fr 1fr' }, gap: 2 }}>
            {filtered.map((item) => {
              const visible = visibleIds.has(item.id);
              const url = item.url.startsWith('http') ? item.url : `https://${item.url}`;
              const groupName = getGroupName(item.group_id);
              return (
                <Card key={item.id} variant="outlined">
                  <CardContent>
                    <Box sx={{ display: 'flex', alignItems: 'center', mb: 1 }}>
                      <Typography variant="h6" sx={{ flexGrow: 1 }} noWrap>{item.title}</Typography>
                      <IconButton size="small" onClick={() => { setEditing(item); setDialogOpen(true); }}>
                        <Edit fontSize="small" />
                      </IconButton>
                      <IconButton size="small" color="error" onClick={() => setDeleteId(item.id)}>
                        <Delete fontSize="small" />
                      </IconButton>
                    </Box>
                    {groupName && (
                      <Chip size="small" label={groupName} variant="outlined" sx={{ mb: 1 }} />
                    )}
                    <Typography variant="body2" color="text.secondary">
                      Логин: <b>{item.username}</b>
                    </Typography>
                    <Box sx={{ display: 'flex', alignItems: 'center', mt: 1, gap: 1 }}>
                      <Typography variant="body2" color="text.secondary">Пароль:</Typography>
                      <Typography component="span" sx={{ fontFamily: 'monospace', fontWeight: 600 }}>
                        {visible ? item.password : '•'.repeat(Math.max(item.password.length, 8))}
                      </Typography>
                      <ToggleVisibilityButton visible={visible} onToggle={() => toggleVisibility(item.id)} />
                      <CopyButton value={item.password} tooltip="Копировать пароль" />
                    </Box>
                    {item.url && (
                      <Box sx={{ mt: 1 }}>
                        <Tooltip title={item.url}>
                          <Chip size="small" label={item.url} variant="outlined"
                            clickable component="a" href={url} target="_blank" rel="noopener" />
                        </Tooltip>
                      </Box>
                    )}
                    {item.notes && (
                      <Typography variant="body2" sx={{ mt: 1, whiteSpace: 'pre-wrap' }} color="text.secondary">
                        {item.notes}
                      </Typography>
                    )}
                  </CardContent>
                  <CardActions sx={{ justifyContent: 'flex-end' }}>
                    <Typography variant="caption" color="text.secondary">
                      Обновлено: {new Date(item.updated_at).toLocaleString()}
                    </Typography>
                  </CardActions>
                </Card>
              );
            })}
          </Box>
        )}
      </Box>

      <Dialog open={dialogOpen} onClose={() => !submitting && setDialogOpen(false)} fullWidth maxWidth="sm">
        <DialogTitle>{editing ? 'Редактировать пароль' : 'Новый пароль'}</DialogTitle>
        <DialogContent>
          <PasswordForm
            initial={editing}
            groups={groups}
            onSubmit={handleCreateOrUpdate}
            onCancel={() => { setDialogOpen(false); setEditing(null); }}
            onRequestCreateGroup={() => setGroupDialog({ open: true, name: '' })}
            submitting={submitting}
          />
        </DialogContent>
      </Dialog>

      <Dialog open={deleteId != null} onClose={() => setDeleteId(null)}>
        <DialogTitle>Удалить пароль?</DialogTitle>
        <DialogContent><DialogContentText>Это действие необратимо.</DialogContentText></DialogContent>
        <DialogActions>
          <Button onClick={() => setDeleteId(null)}>Отмена</Button>
          <Button color="error" variant="contained" onClick={handleDelete}>Удалить</Button>
        </DialogActions>
      </Dialog>

      <Dialog open={groupDialog.open} onClose={() => setGroupDialog({ open: false, name: '' })} fullWidth maxWidth="xs">
        <DialogTitle>Новая группа</DialogTitle>
        <DialogContent>
          <TextField
            autoFocus fullWidth margin="dense" label="Название"
            value={groupDialog.name}
            onChange={(e) => setGroupDialog((p) => ({ ...p, name: e.target.value }))}
            onKeyDown={(e) => { if (e.key === 'Enter') void createGroup(); }}
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setGroupDialog({ open: false, name: '' })}>Отмена</Button>
          <Button variant="contained" onClick={createGroup}>Создать</Button>
        </DialogActions>
      </Dialog>

      <Dialog open={deleteGroupId != null} onClose={() => setDeleteGroupId(null)}>
        <DialogTitle>Удалить группу?</DialogTitle>
        <DialogContent>
          <DialogContentText>Пароли из этой группы не удалятся — они останутся без группы.</DialogContentText>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setDeleteGroupId(null)}>Отмена</Button>
          <Button color="error" variant="contained" onClick={removeGroup}>Удалить</Button>
        </DialogActions>
      </Dialog>

      <Snackbar open={!!snack} autoHideDuration={3000} onClose={() => setSnack(null)} message={snack ?? ''} />
    </Box>
  );
}
