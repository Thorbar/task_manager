<script lang="ts">
  import { t } from 'svelte-i18n';
  import { goto } from '$app/navigation';
  import Button from '$lib/components/Button.svelte';
  import InputWrapper from '$lib/components/InputWrapper.svelte';
  import TextLink from '$lib/components/TextLink.svelte';
  import ButtonLoadSpinner from '$lib/components/Button.svelte';
  import Spinner from '$lib/components/Spinner.svelte';
  import { PUBLIC_API_URL, PUBLIC_API_PORT } from '$env/static/public';

  // User credentials state
  let username = '';
  let password = '';

  // UI state
  let errorMsg = '';
  let loading = false;

  // Base API URL from environment variables
  const baseUrl = `${PUBLIC_API_URL}:${PUBLIC_API_PORT}`;

  /**
   * Handles user login submission.
   * Sends credentials to backend, stores token on success,
   * redirects to tasks page. Shows errors on failure.
   */
  const login = async (e: Event) => {
    e.preventDefault();
    errorMsg = '';
    loading = true;

    try {
      const res = await fetch(`${baseUrl}/login`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ username, password })
      });

      if (!res.ok) {
        // Extract backend error message text
        const err = await res.text();
        errorMsg = err;
        loading = false;
        return;
      }

      const data = await res.json();
      localStorage.setItem('token', data.token);

      // Delay to show loading spinner for UX
      setTimeout(() => {
        loading = false;
        goto('/tasks');
      }, 1000);

    } catch (err: unknown) {
      loading = false;
      if (err instanceof Error) {
        errorMsg = err.message;
      } else {
        errorMsg = $t('unknown');
      }
    }
  };

  /**
   * Navigate to the registration page
   */
  function register() {
    goto('/register');
  }
</script>

<div class="flex flex-col items-center mt-32">
  <div class="text-xl font-semibold mb-6 uppercase text-green-600 drop-shadow-lg">
    <h2>{$t('login')}</h2>
  </div>

  <form on:submit|preventDefault={login} class="space-y-5 w-full flex flex-col items-center">
    <div class="w-55">
      <label for="username" class="block font-bold text-sm mb-1">{$t('username')}</label>
      <InputWrapper>
        <input id="username" bind:value={username} type="text" placeholder={$t('insertUser')} required class="w-full"/>
      </InputWrapper>
    </div>

    <div class="w-55">
      <label for="password" class="block font-bold text-sm mb-1">{$t('password')}</label>
      <InputWrapper>
        <input id="password" bind:value={password} type="password" placeholder={$t('insertPassword')} required class="w-full"/>
      </InputWrapper>
    </div>

    <TextLink><a href="/resetPassword">{$t('resetPassword')}</a></TextLink>

    <div class="h-15 flex justify-between items-center w-55 text-center">
      <ButtonLoadSpinner type="submit" loading={loading}>{$t('login')}</ButtonLoadSpinner>
      <TextLink><a href="/">{ $t('main') }</a></TextLink>
    </div>

    <div class="flex justify-between items-center w-55 text-center">
      <Button on:click={register}>{$t('register')}</Button>
    </div>
  </form>

  {#if errorMsg}
    <p class="bg-red-600 text-white text-center font-bold text-sm whitespace-pre-line px-4 py-2 rounded shadow-lg z-50 transition-opacity duration-300">
      {errorMsg}
    </p>
  {/if}

  {#if loading}
    <Spinner />
  {/if}
</div>
