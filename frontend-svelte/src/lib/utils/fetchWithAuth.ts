// src/lib/utils/fetchWithAuth.ts
import { goto } from '$app/navigation';
import { get } from 'svelte/store';
import { t } from 'svelte-i18n';


export async function fetchWithAuth(input: RequestInfo, init: RequestInit = {}) {
  const token = localStorage.getItem('token');

  const res = await fetch(input, {
    ...init,
    headers: {
      ...init.headers,
      Authorization: token ? `Bearer ${token}` : ''
    }
  });

  if (res.status === 401) {
    localStorage.removeItem('token');
    goto('/');

        // Use get(t) to retrieve current translation string synchronously
    throw new Error(get(t)('token_expired'));
  }

  return res;
}
