'use client';
import { cn } from '@/lib/utils';
import React, { useMemo } from 'react';

function seededRandom(seed: number) {
  const x = Math.sin(seed) * 10000;
  return x - Math.floor(x);
}

export const Meteors = ({
  number,
  className,
}: {
  number?: number;
  className?: string;
}) => {
  const meteors = useMemo(() => {
    const count = number || 20;
    return Array.from({ length: count }, (_, index) => {
      const seed = index + 1; // Use index + 1 as seed to avoid 0
      return {
        id: `meteor-${index}`,
        left: `${Math.floor(seededRandom(seed) * 100)}%`,
        delay: `${seededRandom(seed + 1) * (0.8 - 0.2) + 0.2}s`,
        duration: `${Math.floor(seededRandom(seed + 2) * (10 - 2) + 2)}s`,
      };
    });
  }, [number]);

  return (
    <div
      className={cn(
        'fixed inset-0 overflow-hidden pointer-events-none',
        className
      )}
    >
      {meteors.map((meteor) => (
        <span
          key={meteor.id}
          className={cn(
            'absolute h-0.5 w-0.5 bg-slate-500 rotate-[215deg] animate-meteor',
            'before:content-[""] before:absolute before:top-1/2 before:transform before:-translate-y-[50%] before:w-[50px] before:h-[1px] before:bg-gradient-to-r before:from-[#64748b] before:to-transparent'
          )}
          style={{
            top: '0px',
            left: meteor.left,
            animationDelay: meteor.delay,
            animationDuration: meteor.duration,
          }}
        />
      ))}
    </div>
  );
};
