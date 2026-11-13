const TOKEN_KEY = 'pwm_token';
const USERNAME_KEY = 'pwm_username';
const EXPIRES_KEY = 'pwm_expires';

export const authStorage = {
  setSession(token: string, username: string, expiresAt: string) {
    localStorage.setItem(TOKEN_KEY, token);
    localStorage.setItem(USERNAME_KEY, username);
    localStorage.setItem(EXPIRES_KEY, expiresAt);
  },
  getToken(): string | null {
    const token = localStorage.getItem(TOKEN_KEY);
    if (!token) return null;
    const exp = localStorage.getItem(EXPIRES_KEY);
    if (exp && new Date(exp).getTime() < Date.now()) {
      this.clear();
      return null;
    }
    return token;
  },
  getUsername(): string | null {
    return localStorage.getItem(USERNAME_KEY);
  },
  clear() {
    localStorage.removeItem(TOKEN_KEY);
    localStorage.removeItem(USERNAME_KEY);
    localStorage.removeItem(EXPIRES_KEY);
  },
};
