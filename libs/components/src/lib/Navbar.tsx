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
import { useAuthStore } from '@trio/hooks';
import {
  Github,
  LogOut,
  MessageCircleHeart,
  Plus,
  Settings,
  User,
} from 'lucide-react';

import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuGroup,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuPortal,
  DropdownMenuSeparator,
  DropdownMenuShortcut,
  DropdownMenuSub,
  DropdownMenuSubContent,
  DropdownMenuSubTrigger,
  DropdownMenuTrigger,
} from '@/shadcn/ui/dropdown-menu';

type NavLink = { title: string; icon?: ReactElement; href: string };
const NavLinks: NavLink[] = [
  { title: 'Home', icon: <IoMdHome />, href: '/' },
  { title: 'Research', href: '/research', icon: <SiRoamresearch /> },
  { title: 'About', href: '/about', icon: <CiSquareQuestion /> },
  { title: 'Chat', icon: <BsChatHeart />, href: '/chat' },
];

export function Navbar({ className }: { className?: string }) {
  const user = useAuthStore((state) => state.user);
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
      <ul
        className={`list-none text-gray-200 md:flex items-center gap-5 hidden ${
          !user && 'pl-[170px]'
        }`}
      >
        {NavLinks.map((link) => (
          <li key={v4()}>
            <Link href={link.href} className="flex items-center">
              {link.title}
            </Link>
          </li>
        ))}
      </ul>
      {user ? (
        <Dropdown>
          <div className="w-10 h-10 rounded-full bg-neutral-800 text-white md:flex hidden justify-center items-center">
            {user.fullName.split(' ').map((name) => name.charAt(0))}
          </div>
        </Dropdown>
      ) : (
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
      )}
      <MobileNav>
        <IoIosMenu className="text-2xl" />
      </MobileNav>
    </nav>
  );
}

const MobileNav = ({ children }: { children: ReactNode }) => {
  const pathname = usePathname();
  const user = useAuthStore((state) => state.user);

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
            <ul>
              {user ? (
                <li className="w-full">
                  <Link
                    href={'/signout'}
                    className="flex gap-2 items-center w-full py-2 text-red-500"
                  >
                    <MdLogin />
                    Logout
                  </Link>
                </li>
              ) : (
                <>
                  {' '}
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
                </>
              )}
            </ul>
          </div>
        </SheetFooter>
      </SheetContent>
    </Sheet>
  );
};

function Dropdown({ children }: { children: ReactNode }) {
  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>{children}</DropdownMenuTrigger>
      <DropdownMenuContent className="w-56">
        <DropdownMenuLabel>My Account</DropdownMenuLabel>
        <DropdownMenuSeparator />
        <DropdownMenuGroup>
          <DropdownMenuItem disabled>
            <User className="mr-2 h-4 w-4" />
            <span>Profile</span>
          </DropdownMenuItem>
          <DropdownMenuItem disabled>
            <Settings className="mr-2 h-4 w-4" />
            <span>Settings</span>
          </DropdownMenuItem>
        </DropdownMenuGroup>
        <DropdownMenuSeparator />
        <DropdownMenuGroup>
          <DropdownMenuItem disabled>
            <Plus className="mr-2 h-4 w-4" />
            <span>New Chat</span>
          </DropdownMenuItem>
        </DropdownMenuGroup>
        <DropdownMenuSeparator />
        <DropdownMenuItem>
          <Link
            href={'https://github.com/somtojf/trio'}
            target="_blank"
            className="space-x-2 flex items-center"
          >
            <Github className=" h-4 w-4" />
            <span>GitHub</span>
          </Link>
        </DropdownMenuItem>
        <DropdownMenuItem>
          <Link
            href="mailto:somtochukwujf@gmail.com"
            target="_blank"
            className="space-x-2 flex items-center"
          >
            <MessageCircleHeart className="h-4 w-4" />
            <span>Feedback</span>
          </Link>
        </DropdownMenuItem>
        <DropdownMenuSeparator />
        <DropdownMenuItem>
          <Link
            href="/signout"
            className="space-x-2 flex items-center text-red-500"
          >
            <LogOut className=" h-4 w-4" />
            <span>Log out</span>
          </Link>
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  );
}
