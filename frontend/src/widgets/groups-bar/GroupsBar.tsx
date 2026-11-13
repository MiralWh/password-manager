import { Box, Chip } from '@mui/material';
import { Add } from '@mui/icons-material';
import { Group } from '@/entities/group/model';

interface Props {
  groups: Group[];
  selectedId: number | null;
  onSelect: (id: number | null) => void;
  onCreate: () => void;
  onDelete: (id: number) => void;
}

export function GroupsBar({ groups, selectedId, onSelect, onCreate, onDelete }: Props) {
  return (
    <Box sx={{ display: 'flex', gap: 1, flexWrap: 'wrap', mb: 2 }}>
      <Chip
        label="Все пароли"
        color={selectedId === null ? 'primary' : 'default'}
        onClick={() => onSelect(null)}
      />
      {groups.map((g) => (
        <Chip
          key={g.id}
          label={`${g.name} (${g.password_count})`}
          color={selectedId === g.id ? 'primary' : 'default'}
          onClick={() => onSelect(g.id)}
          onDelete={() => onDelete(g.id)}
        />
      ))}
      <Chip
        icon={<Add />}
        label="Группа"
        variant="outlined"
        onClick={onCreate}
      />
    </Box>
  );
}
