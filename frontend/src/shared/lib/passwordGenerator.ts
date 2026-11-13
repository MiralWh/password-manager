export interface GeneratorOptions {
  length: number;
  uppercase: boolean;
  lowercase: boolean;
  digits: boolean;
  symbols: boolean;
  excludeAmbiguous: boolean;
}

export const defaultGeneratorOptions: GeneratorOptions = {
  length: 16,
  uppercase: true,
  lowercase: true,
  digits: true,
  symbols: true,
  excludeAmbiguous: true,
};

const UPPER     = 'ABCDEFGHIJKLMNOPQRSTUVWXYZ';
const LOWER     = 'abcdefghijklmnopqrstuvwxyz';
const DIGITS    = '0123456789';
const SYMBOLS   = '!@#$%^&*()-_=+[]{};:,.<>/?';
const AMBIGUOUS = 'O0oIl1';

export function generatePassword(opts: GeneratorOptions): string {
  const strip = (s: string) =>
    opts.excludeAmbiguous ? s.split('').filter((c) => !AMBIGUOUS.includes(c)).join('') : s;

  let pool = '';
  if (opts.uppercase) pool += strip(UPPER);
  if (opts.lowercase) pool += strip(LOWER);
  if (opts.digits)    pool += strip(DIGITS);
  if (opts.symbols)   pool += SYMBOLS;
  if (!pool)          pool = LOWER;

  const length = Math.max(1, Math.min(128, opts.length));
  const buf = new Uint32Array(length);
  crypto.getRandomValues(buf);
  return Array.from(buf, (v) => pool[v % pool.length]).join('');
}
