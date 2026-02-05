import type { Metadata } from 'next';
import './globals.css';

export const metadata: Metadata = {
  title: 'Drive Course Viewer',
  description: 'View Google Drive videos as courses',
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="vi">
      <body className="antialiased min-h-screen">{children}</body>
    </html>
  );
}
