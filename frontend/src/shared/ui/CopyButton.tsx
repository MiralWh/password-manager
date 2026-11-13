import { useState } from 'react';
import { IconButton, Tooltip } from '@mui/material';
import { Check, ContentCopy } from '@mui/icons-material';
import { copyToClipboard } from '../lib/clipboard';

interface Props {
  value: string;
  size?: 'small' | 'medium';
  tooltip?: string;
}

export function CopyButton({ value, size = 'small', tooltip = 'Копировать' }: Props) {
  const [copied, setCopied] = useState(false);

  const onClick = async () => {
    const ok = await copyToClipboard(value);
    if (!ok) return;
    setCopied(true);
    setTimeout(() => setCopied(false), 1500);
  };

  return (
    <Tooltip title={copied ? 'Скопировано!' : tooltip}>
      <IconButton size={size} onClick={onClick} type="button">
        {copied ? <Check fontSize="small" color="success" /> : <ContentCopy fontSize="small" />}
      </IconButton>
    </Tooltip>
  );
}
