<script lang="ts">
  import { t } from 'svelte-i18n';
  import { goto } from '$app/navigation';
  import { PUBLIC_API_URL, PUBLIC_API_PORT } from '$env/static/public';
  import InputWrapper from '$lib/components/InputWrapper.svelte';
  import Button from '$lib/components/Button.svelte';
  import Spinner from '$lib/components/Spinner.svelte';

  // State variables
  let email = '';
  let errorMsg = '';
  let showError = false;
  let successMsg = '';    
  let showSuccess = false;
  let loading = false;

  const baseUrl = `${PUBLIC_API_URL}:${PUBLIC_API_PORT}`;

  /**
   * Displays an error message for 2 seconds.
   * @param msg The message to display
   */
  function showErrorMessage(msg: string) {
    errorMsg = msg;
    showError = true;
    setTimeout(() => showError = false, 2000);
  }

  /**
   * Displays a success message for 2 seconds.
   * @param msg The message to display
   */
  function showSuccessMessage(msg: string) {
    successMsg = msg;
    showSuccess = true;
    setTimeout(() => showSuccess = false, 2000);
  }

  /**
   * Handles password reset confirmation request.
   * Sends the email to the backend and handles success or error responses.
   */
  async function handleConfirm() {
    loading = true;

    try {
      const res = await fetch(`${baseUrl}/forgot-password`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email })
      });

      if (!res.ok) {
        const data = await res.json();
        // Use backend error message if available, fallback to translation
        errorMsg = data.error || $t('resetPwdFailed');
        showErrorMessage(errorMsg);
        loading = false;
        return;
      }

      showSuccessMessage($t('resetPwdSuccess'));
      setTimeout(() => {
        loading = false;
        goto('/login');
      }, 1000);

    } catch {
      showErrorMessage($t('serverError'));
      loading = false;
    }
  }

  /** Redirects to login page */
  function handleCancel() {
    goto('/login');
  }
</script>

<div class="flex flex-col items-center mt-32">
  <h2 class="text-xl font-semibold mb-6 uppercase text-green-600 drop-shadow-lg">{$t('pwdRecovery')}</h2>

  <div class="h-14 w-55">
    <label for="email" class="block text-sm font-bold mb-1">{$t('email')}</label>
    <InputWrapper>
      <input id="email" type="email" bind:value={email} required class="w-full" />
    </InputWrapper>
  </div>

  {#if showError}
    <p class="text-red-600 font-bold text-sm mt-2">{errorMsg}</p>
  {/if}

  {#if showSuccess}
    <p class="text-green-600 font-bold text-sm mt-2">{successMsg}</p>
  {/if}

  <div class="flex justify-between items-center gap-5 w-55 mt-5 text-center">
    <Button on:click={handleConfirm}>{$t('confirm')}</Button>
    <Button on:click={handleCancel} variant="secondary">{$t('cancel')}</Button>
  </div>

  {#if loading}
    <Spinner />
  {/if}
</div>
