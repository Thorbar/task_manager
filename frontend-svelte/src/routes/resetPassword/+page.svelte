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

<section class="flex flex-col items-center">
  <form on:submit|preventDefault={handleConfirm} class="space-y-3 w-full flex flex-col items-center">

    <h2 class="text-3xl font-extrabold uppercase tracking-wide text-green-600 drop-shadow-lg mb-10">{$t('pwdRecovery')}</h2>

    <div class="">
      <label for="email" class="block text-sm font-bold mb-1">{$t('email')}</label>
      <InputWrapper>
        <input id="email" type="email" bind:value={email}  required class="w-full" />
      </InputWrapper>
    </div>


    <div class="flex justify-between items-center gap-5 w-55 mt-5 text-center">
      <Button on:click={handleConfirm}>{$t('confirm')}</Button>
      <Button on:click={handleCancel} variant="secondary">{$t('cancel')}</Button>
    </div>


  </form>
    {#if loading}
      <Spinner />
    {/if}
    {#if showError}
      <p class="bg-red-600 text-white text-center font-bold text-sm whitespace-pre-line px-4 py-2 rounded shadow-lg z-50 transition-opacity duration-300">{errorMsg}</p>
    {/if}

    {#if showSuccess}
      <p class="bg-green-600 text-black text-center font-bold text-sm whitespace-pre-line px-4 py-2 rounded shadow-lg z-50 transition-opacity duration-300">{successMsg}</p>
    {/if}
</section>

