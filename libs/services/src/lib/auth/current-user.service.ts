import { BaseRoute, Route } from '../routes';
import { User } from '@trio/types';

export const getCurrentUser = async (): Promise<User> => {
  const res = await fetch(`${BaseRoute}/${Route.Me}`, {
    method: 'GET',
    credentials: 'include',
  });
  if (!res.ok) throw new Error(res.statusText);
  const user = await res.json();
  return user.data as User;
};
