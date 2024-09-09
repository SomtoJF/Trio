import SpotlightWrapper from '../components/SpotlightWrapper';
import './global.css';
import { Inter } from 'next/font/google';
import Providers from './providers';

const inter = Inter({ subsets: ['latin'] });

export const metadata = {
  title: 'Welcome to trio',
  description:
    'Trio is a multi-agent chat app which provides you with multiple AI friends to help you with your questions',
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en" className={inter.className}>
      <Providers>
        <body className="w-screen bg-black text-white flex justify-center">
          <SpotlightWrapper />
          <div className="max-w-screen-2xl w-full">{children}</div>
        </body>
      </Providers>
    </html>
  );
}
