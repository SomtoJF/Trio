'use client';

import { Button } from '@/shadcn/ui/button';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import React from 'react';
import { IoMdRefresh } from 'react-icons/io';

export default function Index() {
  const router = useRouter();
  return (
    <div className="w-full h-full flex justify-center items-center">
      <div className="space-y-3">
        <p className="font-bold text-2xl text-center">An error occured</p>
        <div className="flex justify-center items-center gap-2">
          <Button
            className="bg-white text-black hover:bg-white active:bg-gray-600"
            onClick={() => router.refresh()}
          >
            Refresh <IoMdRefresh className="ml-1" />
          </Button>
          <Button className="bg-white text-black hover:bg-white">
            <Link href={'/'} className="w-full h-full">
              Go Home
            </Link>
          </Button>
        </div>
      </div>
    </div>
  );
}
