import { writable } from "svelte/store";

export interface ToastData {
  header?: string;
  body: string;
  duration?: number;
}

interface IdToastData extends ToastData {
  id: number;
}

export const toastStore = writable<IdToastData[]>([]);

let maxId = 0;

export const addToast = (data: ToastData) => {
  const id = maxId++;

  console.log(id);

  toastStore.update((toasts) => (toasts = [...toasts, { ...data, id }]));

  if (data.duration !== -1) {
    setTimeout(() => {
      toastStore.update((toasts) => (toasts = toasts.filter((toast) => toast.id !== id)));
    }, data.duration ?? 5000);
  }

  return id;
};

export const delToast = (id: number) => {
  toastStore.update((toasts) => (toasts = toasts.filter((toast) => toast.id !== id)));
};
