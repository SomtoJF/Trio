import { BaseRoute, Route } from '../routes';

export const signout = async () => {
  const res = await fetch(`${BaseRoute}/${Route.SignOut}`, {
    method: 'POST',
    credentials: 'include',
  });
  if (res.status > 299) throw new Error(res.statusText);
};
