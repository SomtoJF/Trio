'use client';

import { Button } from '@/shadcn/ui/button';
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/shadcn/ui/form';
import { Input } from '@/shadcn/ui/input';
import { zodResolver } from '@hookform/resolvers/zod';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { AiOutlineLoading } from 'react-icons/ai';
import Link from 'next/link';
import { useForm } from 'react-hook-form';
import { z } from 'zod';
import { login } from '@trio/services';
import { queryKeys } from '@trio/query-key-factory';
import { useToast } from '@trio/hooks';

const formSchema = z.object({
  userName: z
    .string()
    .min(2, {
      message: 'Username must be at least 2 characters.',
    })
    .max(20, { message: 'username cannot be more than 20 characters' }),
  password: z
    .string()
    .min(8, {
      message: 'Password must be at least 8 characters.',
    })
    .max(20, { message: 'Password cannot be more than 20 characters' }),
});

export default function Page() {
  const toast = useToast();
  const queryClient = useQueryClient();
  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      userName: '',
      password: '',
    },
  });

  const loginMutation = useMutation({
    mutationFn: (data: z.infer<typeof formSchema>) => login(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.user.all });
      toast.success('Login successful');
    },
    onError: (error) => {
      toast.error(error.message);
      throw error;
    },
  });

  function onSubmit(values: z.infer<typeof formSchema>) {
    loginMutation.mutate(values);
  }

  return (
    <div className="w-full h-full flex justify-end">
      <section className="w-full p-10 lg:p-20 space-y-8 bg-white text-black rounded-xl lg:h-full">
        <div>
          <h2 className="font-bold text-2xl">Welcome back</h2>
          <p className="text-gray-500">
            Ready to continue the conversation? Log in and pick up right where
            you left off.
          </p>
        </div>
        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-8">
            <FormField
              control={form.control}
              name="userName"
              render={({ field }) => (
                <FormItem>
                  <FormLabel className="font-bold text-sm">Username</FormLabel>
                  <FormControl>
                    <Input placeholder="Your username..." {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name="password"
              render={({ field }) => (
                <FormItem>
                  <FormLabel className="font-bold text-sm">Password</FormLabel>
                  <FormControl>
                    <Input placeholder="*******" {...field} type="password" />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <small>
              Dont have an account?{' '}
              <Link href="/signup" className="font-bold">
                Sign up
              </Link>
            </small>
            <Button
              type="submit"
              className="w-full"
              disabled={loginMutation.isPending}
            >
              {loginMutation.isPending ? (
                <AiOutlineLoading className="animate-spin" />
              ) : (
                'Submit'
              )}
            </Button>
          </form>
        </Form>
      </section>
    </div>
  );
}
