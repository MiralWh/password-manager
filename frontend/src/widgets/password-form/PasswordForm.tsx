import { ChangeEvent, FormEvent, useEffect, useState } from 'react';
import {
  Box, Button, Divider, FormControl, IconButton, InputAdornment, InputLabel,
  MenuItem, Select, TextField, Tooltip,
} from '@mui/material';
import { Casino, Visibility, VisibilityOff } from '@mui/icons-material';
import { PasswordEntry, PasswordInput } from '@/entities/password/model';
import { Group } from '@/entities/group/model';
import { CopyButton } from '@/shared/ui/CopyButton';
import { GeneratePasswordDialog } from '@/features/passwords/generate-password/GeneratePasswordDialog';

interface Props {
  initial?: PasswordEntry | null;
  groups: Group[];
  onSubmit: (input: PasswordInput) => Promise<void>;
  onCancel: () => void;
  onRequestCreateGroup: () => void;
  submitting?: boolean;
}

const empty: PasswordInput = {
  group_id: null, title: '', username: '', password: '', url: '', notes: '',
};

export function PasswordForm({
  initial, groups, onSubmit, onCancel, onRequestCreateGroup, submitting,
}: Props) {
  const [form, setForm] = useState<PasswordInput>(empty);
  const [show, setShow] = useState(false);
  const [generatorOpen, setGeneratorOpen] = useState(false);

  useEffect(() => {
    if (initial) {
      setForm({
        group_id: initial.group_id ?? null,
        title: initial.title,
        username: initial.username,
        password: initial.password,
        url: initial.url,
        notes: initial.notes,
      });
    } else {
      setForm(empty);
    }
    setShow(false);
  }, [initial]);

  const onChange = (field: keyof PasswordInput) => (e: ChangeEvent<HTMLInputElement>) =>
    setForm((p) => ({ ...p, [field]: e.target.value }));

  const submit = async (e: FormEvent) => {
    e.preventDefault();
    await onSubmit(form);
  };

  return (
    <>
      <Box component="form" onSubmit={submit} sx={{ display: 'flex', flexDirection: 'column', gap: 2, mt: 1 }}>
        <FormControl fullWidth>
          <InputLabel>Группа</InputLabel>
          <Select
            label="Группа"
            value={form.group_id ?? ''}
            onChange={(e) => {
              const v = e.target.value as string | number;
              if (v === '__create__') { onRequestCreateGroup(); return; }
              setForm((p) => ({ ...p, group_id: v === '' ? null : Number(v) }));
            }}
          >
            <MenuItem value="">Без группы</MenuItem>
            {groups.map((g) => (
              <MenuItem key={g.id} value={g.id}>{g.name}</MenuItem>
            ))}
            <Divider />
            <MenuItem value="__create__"><em>+ Создать группу</em></MenuItem>
          </Select>
        </FormControl>

        <TextField label="Название" value={form.title} onChange={onChange('title')} required fullWidth />
        <TextField label="Логин/Пользователь" value={form.username} onChange={onChange('username')} required fullWidth />

        <TextField
          label="Пароль"
          type={show ? 'text' : 'password'}
          value={form.password}
          onChange={onChange('password')}
          required
          fullWidth
          InputProps={{
            endAdornment: (
              <InputAdornment position="end">
                <Tooltip title="Сгенерировать пароль">
                  <IconButton type="button" onClick={() => setGeneratorOpen(true)}>
                    <Casino />
                  </IconButton>
                </Tooltip>
                <CopyButton value={form.password} />
                <IconButton type="button" onClick={() => setShow((v) => !v)}>
                  {show ? <VisibilityOff /> : <Visibility />}
                </IconButton>
              </InputAdornment>
            ),
          }}
        />

        <TextField label="URL" value={form.url} onChange={onChange('url')} fullWidth />
        <TextField label="Заметки" value={form.notes} onChange={onChange('notes')} fullWidth multiline minRows={2} />

        <Box sx={{ display: 'flex', gap: 1, justifyContent: 'flex-end' }}>
          <Button onClick={onCancel} disabled={submitting}>Отмена</Button>
          <Button type="submit" variant="contained" disabled={submitting}>
            {initial ? 'Сохранить' : 'Создать'}
          </Button>
        </Box>
      </Box>

      <GeneratePasswordDialog
        open={generatorOpen}
        onClose={() => setGeneratorOpen(false)}
        onUse={(pwd) => setForm((p) => ({ ...p, password: pwd }))}
      />
    </>
  );
}
