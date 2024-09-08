import SpotlightWrapper from '../components/SpotlightWrapper';
import './global.css';

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
    <html lang="en">
      <body className="w-screen bg-black flex justify-center">
        <SpotlightWrapper />
        <div className="max-w-screen-2xl">{children}</div>
      </body>
    </html>
  );
}
