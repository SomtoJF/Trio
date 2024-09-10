import { create } from 'zustand';
import { User } from '@trio/types';

interface UserState {
  user: User | null;
  setUser: (user: User | null) => void;
}

export const useAuthStore = create<UserState>()((set) => ({
  user: null,
  setUser: (user) => set(() => ({ user: user })),
}));
