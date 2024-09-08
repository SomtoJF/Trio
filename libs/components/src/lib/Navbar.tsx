import Link from 'next/link';

export function Navbar() {
  return (
    <nav className="w-full h-16 flex justify-between items-center text-sm relative z-20">
      <h1 className="font-extrabold text-2xl">Trio</h1>
      <ul className="list-none text-gray-200 flex items-center gap-5 pl-[170px]">
        <li>Research</li>
        <li>About</li>
        <li>Chat</li>
      </ul>
      <div className="flex gap-2 items-center">
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
    </nav>
  );
}
