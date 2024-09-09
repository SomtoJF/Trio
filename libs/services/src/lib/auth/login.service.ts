import { BaseRoute, Route } from '../routes';

export async function login(data: {
  userName: string;
  password: string;
}): Promise<void> {
  const response = await fetch(`${BaseRoute}/${Route.Login}`, {
    method: 'POST',
    credentials: 'include',
    body: JSON.stringify(data),
  });
  if (response.status > 299) throw new Error(response.statusText);
}
