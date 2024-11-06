'use client';

import React, { useEffect } from 'react';
import { LoadingScreen } from '@trio/components';
import { useQuery } from '@tanstack/react-query';
import { queryKeys } from '@trio/query-key-factory';
import { getCurrentUser } from '@trio/services';
import { useRouter } from 'next/navigation';
import { useAuthStore } from '@trio/hooks';

export default function Layout({ children }: { children: React.ReactNode }) {
  const { push } = useRouter();
  const setUser = useAuthStore((state) => state.setUser);
  const currentUserQuery = useQuery({
    queryKey: queryKeys.user.currentUser(),
    queryFn: getCurrentUser,
    retry: false,
  });

  useEffect(() => {
    if (!currentUserQuery.isFetching) setUser(currentUserQuery.data ?? null);
  }, [currentUserQuery.isFetching, currentUserQuery.data]);

  if (currentUserQuery.isPending || currentUserQuery.isFetching) {
    return <LoadingScreen />;
  } else if (currentUserQuery.isError) {
    push('/login');
  } else return children;
}
