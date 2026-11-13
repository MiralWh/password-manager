import { useEffect, useState } from 'react';
import {
  Box, Button, Dialog, DialogActions, DialogContent, DialogTitle,
  FormControlLabel, IconButton, InputAdornment, Slider, Switch, TextField, Tooltip, Typography,
} from '@mui/material';
import { Refresh } from '@mui/icons-material';
import { CopyButton } from '@/shared/ui/CopyButton';
import {
  defaultGeneratorOptions, generatePassword, GeneratorOptions,
} from '@/shared/lib/passwordGenerator';

interface Props {
  open: boolean;
  onClose: () => void;
  onUse: (password: string) => void;
}

export function GeneratePasswordDialog({ open, onClose, onUse }: Props) {
  const [opts, setOpts] = useState<GeneratorOptions>(defaultGeneratorOptions);
  const [preview, setPreview] = useState('');

  useEffect(() => {
    if (open) setPreview(generatePassword(opts));
  }, [open, opts]);

  const regenerate = () => setPreview(generatePassword(opts));
  const setOpt = <K extends keyof GeneratorOptions>(k: K, v: GeneratorOptions[K]) =>
    setOpts((p) => ({ ...p, [k]: v }));

  return (
    <Dialog open={open} onClose={onClose} fullWidth maxWidth="sm">
      <DialogTitle>Генератор пароля</DialogTitle>
      <DialogContent>
        <TextField
          value={preview}
          fullWidth
          margin="dense"
          InputProps={{
            readOnly: true,
            sx: { fontFamily: 'monospace', fontSize: 16 },
            endAdornment: (
              <InputAdornment position="end">
                <Tooltip title="Ещё"><IconButton onClick={regenerate}><Refresh /></IconButton></Tooltip>
                <CopyButton value={preview} tooltip="Копировать пароль" />
              </InputAdornment>
            ),
          }}
        />
        <Box sx={{ mt: 2 }}>
          <Typography gutterBottom>Длина: {opts.length}</Typography>
          <Slider
            value={opts.length}
            min={8}
            max={64}
            step={1}
            valueLabelDisplay="auto"
            onChange={(_, v) => setOpt('length', v as number)}
          />
        </Box>
        <Box sx={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: 1 }}>
          <FormControlLabel control={<Switch checked={opts.uppercase} onChange={(e) => setOpt('uppercase', e.target.checked)} />} label="A-Z" />
          <FormControlLabel control={<Switch checked={opts.lowercase} onChange={(e) => setOpt('lowercase', e.target.checked)} />} label="a-z" />
          <FormControlLabel control={<Switch checked={opts.digits} onChange={(e) => setOpt('digits', e.target.checked)} />} label="0-9" />
          <FormControlLabel control={<Switch checked={opts.symbols} onChange={(e) => setOpt('symbols', e.target.checked)} />} label="!@#$" />
          <FormControlLabel control={<Switch checked={opts.excludeAmbiguous} onChange={(e) => setOpt('excludeAmbiguous', e.target.checked)} />} label="Исключить O0oIl1" />
        </Box>
      </DialogContent>
      <DialogActions>
        <Button onClick={onClose}>Отмена</Button>
        <Button variant="contained" onClick={() => { onUse(preview); onClose(); }}>
          Использовать
        </Button>
      </DialogActions>
    </Dialog>
  );
}
