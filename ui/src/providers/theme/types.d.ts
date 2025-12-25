export type ThemeProviderProps = {
  children: React.ReactNode;
  defaultTheme?: Theme;
  storageKey?: string;
};

export type ThemeProviderProps = {
  children: React.ReactNode;
  defaultTheme?: Theme;
  storageKey?: string;
};

export type Theme = "dark" | "light" | "system";

export type ThemeProviderState = {
  theme: Theme;
  setTheme: (theme: Theme) => void;
};
