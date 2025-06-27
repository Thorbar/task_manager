// src/lib/stores/tasks.ts
import { writable } from 'svelte/store';
import type { Task } from '$lib/types';

export const tasks = writable<Task[]>([]);
export const errorMsg = writable('');
