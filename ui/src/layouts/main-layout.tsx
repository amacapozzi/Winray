import { ThemeProvider } from "@/providers/theme/theme-provider";

interface MainLayoutProps {
  children: React.ReactNode;
}

export const MainLayout = ({ children }: Readonly<MainLayoutProps>) => {
  return (
    <ThemeProvider defaultTheme="dark">
      <div className="w-screen absolute top-0 left-0 right-0 w-full">
        {children}
      </div>
    </ThemeProvider>
  );
};
