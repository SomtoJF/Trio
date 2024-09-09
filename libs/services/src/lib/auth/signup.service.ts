import { BaseRoute, Route } from '../routes';

type SignUpType = {
  userName: string;
  password: string;
  fullName: string;
};

export async function signUp(data: SignUpType): Promise<void> {
  const response = await fetch(`${BaseRoute}/${Route.SignUp}`, {
    method: 'POST',
    credentials: 'include',
    body: JSON.stringify(data),
  });
  if (response.status > 299) throw new Error(response.statusText);
}
