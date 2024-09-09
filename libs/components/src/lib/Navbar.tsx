'use client';

import Link from 'next/link';
import { cn } from '@/lib/utils';
import { IoIosMenu } from 'react-icons/io';
import { BsChatHeart } from 'react-icons/bs';
import { SiRoamresearch } from 'react-icons/si';
import { CiSquareQuestion } from 'react-icons/ci';
import { IoMdHome } from 'react-icons/io';
import { MdLogin } from 'react-icons/md';
import { IoCreateOutline } from 'react-icons/io5';
import {
  Sheet,
  SheetContent,
  SheetFooter,
  SheetHeader,
  SheetTitle,
  SheetTrigger,
} from '@/shadcn/ui/sheet';
import { ReactElement, ReactNode } from 'react';
import { usePathname } from 'next/navigation';
import { v4 } from 'uuid';

type NavLink = { title: string; icon?: ReactElement; href: string };
const NavLinks: NavLink[] = [
  { title: 'Home', icon: <IoMdHome />, href: '/' },
  { title: 'Research', href: '/research', icon: <SiRoamresearch /> },
  { title: 'About', href: '/about', icon: <CiSquareQuestion /> },
  { title: 'Chat', icon: <BsChatHeart />, href: '/chat' },
];

export function Navbar({ className }: { className?: string }) {
  return (
    <nav
      className={cn(
        'w-full h-16 flex justify-between items-center text-sm relative z-20',
        className
      )}
    >
      <Link href={'/'} className="font-extrabold text-2xl">
        Trio
      </Link>
      <ul className="list-none text-gray-200 md:flex items-center gap-5 pl-[170px] hidden">
        {NavLinks.map((link) => (
          <li key={v4()}>
            <Link href={link.href} className="flex items-center">
              {link.title}
            </Link>
          </li>
        ))}
      </ul>
      <div className="gap-2 items-center hidden md:flex">
        <Link
          href="/login"
          className="bg-gray-200 px-8 py-2 inline-flex items-center h-max rounded-lg text-black hover:bg-gray-200"
        >
          Login
        </Link>
        <Link href="/signup" className=" px-8 py-4 bg-transparent">
          Sign Up
        </Link>
      </div>
      <MobileNav>
        <IoIosMenu className="text-2xl" />
      </MobileNav>
    </nav>
  );
}

const MobileNav = ({ children }: { children: ReactNode }) => {
  const pathname = usePathname();

  return (
    <Sheet>
      <SheetTrigger className="md:hidden">{children}</SheetTrigger>
      <SheetContent className="bg-neutral-950 text-white space-y-8 border-none">
        <SheetHeader className="text-gray-200 mb-10">
          <SheetTitle className="text-left text-white">
            <Link href={'/'} className="font-extrabold text-2xl px-3">
              Trio
            </Link>
          </SheetTitle>
        </SheetHeader>

        <div className="space-y-4">
          <h3 className="text-xs font-bold text-gray-500 px-3">MENU</h3>
          <ul className="w-full space-y-2 text-xs font-semibold">
            {NavLinks.map((link) => (
              <li key={link.href}>
                <Link
                  href={link.href}
                  className={`flex gap-2 items-center w-full py-2 px-3 rounded-lg ${
                    pathname === link.href
                      ? 'bg-neutral-800'
                      : 'hover:bg-neutral-800'
                  }`}
                >
                  {link.icon}
                  {link.title}
                </Link>
              </li>
            ))}
          </ul>
        </div>

        <SheetFooter className="absolute bottom-5">
          <div className="text-sm font-bold px-3">
            <h3 className="text-xs font-bold text-gray-500 mb-2">
              YOUR ACCOUNT
            </h3>
            <ul className="">
              <li>
                <Link
                  href={'/login'}
                  className="flex gap-2 items-center w-full py-2"
                >
                  <MdLogin />
                  Login
                </Link>
              </li>
              <li>
                <Link
                  href={'/signup'}
                  className="flex gap-2 items-center w-full py-2"
                >
                  <IoCreateOutline />
                  Signup
                </Link>
              </li>
            </ul>
          </div>
        </SheetFooter>
      </SheetContent>
    </Sheet>
  );
};
