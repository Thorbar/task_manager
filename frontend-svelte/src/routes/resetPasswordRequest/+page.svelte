<script lang="ts">
  import { t } from 'svelte-i18n';
  import { onMount } from 'svelte';
  import { PUBLIC_API_URL, PUBLIC_API_PORT } from '$env/static/public';
  import { goto } from '$app/navigation';
  import InputWrapper from '$lib/components/InputWrapper.svelte';
  import TextLink from '$lib/components/TextLink.svelte';
  import ButtonLoadSpinner from '$lib/components/Button.svelte';

  let token = '';
  let password = '';
  let confirmPassword = '';
  let loading = false;
  let errorMsg = '';
  let showError = false;
  let successMsg = '';
  let showSuccess = false;

  const baseUrl = `${PUBLIC_API_URL}:${PUBLIC_API_PORT}`;

  // Extract token from URL query params on component mount
  onMount(() => {
    const urlParams = new URLSearchParams(window.location.search);
    token = urlParams.get('token') || '';
  });

  /**
   * Validates password strength.
   * Must have at least 8 chars, upper, lower, digit, special char.
   */
  function isValidPassword(pwd: string): boolean {
    const regex = /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[\W_]).{8,}$/;
    return regex.test(pwd);
  }

  /**
   * Shows error message for 2 seconds
   */
  function showErrorMessage(msg: string) {
    errorMsg = msg;
    showError = true;
    setTimeout(() => showError = false, 2000);
  }

  /**
   * Shows success message for 2 seconds
   */
  function showSuccessMessage(msg: string) {
    successMsg = msg;
    showSuccess = true;
    setTimeout(() => showSuccess = false, 2000);
  }

  /**
   * Handles password reset submission.
   * Validates inputs and sends request to backend.
   * Navigates to login on success.
   */
  async function resetPassword() {
    if (!isValidPassword(password)) {
      showErrorMessage($t('wrongPassword'));
      return;
    }

    if (password !== confirmPassword) {
      showErrorMessage($t('passwordMismatch'));
      return;
    }

    loading = true;

    try {
      const res = await fetch(`${baseUrl}/reset-password`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ token, password }),
      });

      const data = await res.json();

      if (!res.ok) {
        loading = false;
        showErrorMessage($t('pwdError'));
      } else {
        showSuccessMessage($t('updatedPassword'));
        setTimeout(() => {
          loading = false;
          goto('/login');
        }, 2000);
      }
    } catch (err) {
      loading = false;
      showErrorMessage($t('serverError'));
      console.error('Error in resetPassword:', err);
    }
  }
</script>

<div class="flex flex-col items-center mt-32">
  <h2 class="text-xl font-semibold mb-6 uppercase text-green-600 drop-shadow-lg">{$t('restablishPwd')}</h2>

  <form on:submit|preventDefault={resetPassword} class="w-55">
    <label for="password" class="block font-bold text-sm mb-1">{$t('newPwd')}</label>
    <InputWrapper>
      <input
        id="password"
        bind:value={password}
        type="password"
        placeholder={$t('insertPassword')}
        required
        class="w-full"
      />
    </InputWrapper>

    <label for="confirmPassword" class="block font-bold text-sm mt-4 mb-1">{$t('confirmPassword')}</label>
    <InputWrapper>
      <input
        id="confirmPassword"
        bind:value={confirmPassword}
        type="password"
        required
        class="w-full"
      />
    </InputWrapper>

    <div class="flex justify-between text-center mt-6">
      <ButtonLoadSpinner type="submit" loading={loading}>{$t('confirm')}</ButtonLoadSpinner>
      <TextLink><a href="/login">{$t('cancel')}</a></TextLink>
    </div>
  </form>

  {#if showError}
    <div
      class="bg-red-600 text-white text-center font-bold text-sm whitespace-pre-line px-4 py-2 rounded shadow-lg z-50 transition-opacity duration-300 mt-4"
    >
      {errorMsg}
    </div>
  {/if}

  {#if showSuccess}
    <div
      class="bg-green-600 text-black text-center font-bold text-sm whitespace-pre-line px-4 py-2 rounded shadow-lg z-50 transition-opacity duration-300 mt-4"
    >
      {successMsg}
    </div>
  {/if}
</div>
