import { IconButton, Tooltip } from '@mui/material';
import { Visibility, VisibilityOff } from '@mui/icons-material';

interface Props { visible: boolean; onToggle: () => void; }

export function ToggleVisibilityButton({ visible, onToggle }: Props) {
  return (
    <Tooltip title={visible ? 'Скрыть пароль' : 'Показать пароль'}>
      <IconButton size="small" onClick={onToggle}>
        {visible ? <VisibilityOff fontSize="small" /> : <Visibility fontSize="small" />}
      </IconButton>
    </Tooltip>
  );
}
