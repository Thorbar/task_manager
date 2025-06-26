<script lang="ts">
  import { t } from 'svelte-i18n';
  import { goto } from '$app/navigation';
  import TextLink from '$lib/components/TextLink.svelte';
  import InputWrapper from '$lib/components/InputWrapper.svelte';
  import ButtonLoadSpinner from '$lib/components/ButtonLoadSpinner.svelte';
  import Spinner from '$lib/components/Spinner.svelte';
  import { PUBLIC_API_URL, PUBLIC_API_PORT } from '$env/static/public';

  // Form fields
  let username = '';
  let password = '';
  let confirmPassword = '';
  let name = '';
  let surname = '';
  let email = '';

  // UI state variables
  let errorMsg = '';
  let showError = false;
  let successMsg = '';
  let showSuccess = false;
  let loading = false;

  const baseUrl = `${PUBLIC_API_URL}:${PUBLIC_API_PORT}`;

  /**
   * Validates password complexity.
   * Password must contain lowercase, uppercase, digit, special character, min length 8.
   */
  function isValidPassword(pwd: string): boolean {
    const regex = /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[\W_]).{8,}$/;
    return regex.test(pwd);
  }

  /**
   * Validates email format and domain against a whitelist of common domains.
   */
  function isValidEmail(email: string): boolean {
    const commonDomains = [
      'gmail.com', 'hotmail.com', 'yahoo.com', 'outlook.com', 'protonmail.com',
      'icloud.com', 'live.com', 'aol.com', 'zoho.com', 'yandex.com', 'mail.ru',
      'gmx.com', 'msn.com', 'me.com', 'ya.ru', 'docomo.ne.jp', 'softbank.ne.jp',
      'rakuten.ne.jp', 'nifty.com', 'biglobe.ne.jp'
    ];

    const basicEmailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    if (!basicEmailRegex.test(email)) return false;

    const domain = email.split('@')[1].toLowerCase();
    return commonDomains.includes(domain);
  }

  /**
   * Displays an error message temporarily.
   */
  function showErrorMessage(msg: string) {
    errorMsg = msg;
    showError = true;
    setTimeout(() => {
      showError = false;
    }, 2000);
  }

  /**
   * Displays a success message temporarily.
   */
  function showSuccessMessage(msg: string) {
    successMsg = msg;
    showSuccess = true;
    setTimeout(() => {
      showSuccess = false;
    }, 2000);
  }

  /**
   * Handles user registration submit event.
   * Validates form and sends POST request to backend.
   */
  const registerUser = async (e: Event) => {
    e.preventDefault();

    if (!isValidPassword(password)) {
      showErrorMessage($t('wrongPassword'));
      return;
    }

    if (password !== confirmPassword) {
      showErrorMessage($t('passwordMismatch'));
      return;
    }

    if (!isValidEmail(email)) {
      showErrorMessage($t('invalidEmail'));
      return;
    }

    try {
      loading = true;

      const res = await fetch(`${baseUrl}/register`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ username, password, name, surname, email })
      });

      if (!res.ok) {
        showErrorMessage($t('registerFailed'));
        loading = false;
        return;
      }

      // Simulate spinner wait time for better UX
      await new Promise(resolve => setTimeout(resolve, 2000));

      loading = false;
      showSuccessMessage($t('registerSuccess'));

      setTimeout(() => {
        goto('/login');
      }, 2000);

    } catch (error) {
      loading = false;
      showErrorMessage($t('registerFailed'));
    }
  };
</script>

<div class="flex flex-col items-center mt-32">
  <div class="text-xl font-semibold mb-6 uppercase text-green-600 drop-shadow-lg">
    <h2>{$t('register')}</h2>
  </div>

  <form on:submit|preventDefault={registerUser}>
    <div class="w-55 font-bold">
      <label for="name" class="block text-sm mb-1">{$t('name')}</label>
      <InputWrapper>
        <input id="name" bind:value={name} type="text" required class="w-full" />
      </InputWrapper>
    </div>

    <div class="h-20 w-55 font-bold">
      <label for="surname" class="block text-sm mb-1">{$t('surname')}</label>
      <InputWrapper>
        <input id="surname" bind:value={surname} type="text" required class="w-full" />
      </InputWrapper>
    </div>

    <div class="w-55 font-bold">
      <label for="email" class="block text-sm mb-1">{$t('email')}</label>
      <InputWrapper>
        <input id="email" bind:value={email} type="email" required class="w-full" />
      </InputWrapper>
    </div>

    <div class="w-55 font-bold">
      <label for="username" class="block text-sm mb-1">{$t('username')}</label>
      <InputWrapper>
        <input id="username" bind:value={username} type="text" required class="w-full" />
      </InputWrapper>
    </div>

    <div class="w-55 font-bold">
      <label for="password" class="block text-sm mb-1">{$t('password')}</label>
      <InputWrapper>
        <input id="password" bind:value={password} type="password" required class="w-full" />
      </InputWrapper>
    </div>

    <div class="h-15 w-55 font-bold">
      <label for="confirmPassword" class="block text-sm mb-1">{$t('confirmPassword')}</label>
      <InputWrapper>
        <input id="confirmPassword" bind:value={confirmPassword} type="password" required class="w-full" />
      </InputWrapper>
    </div>

    <div class="flex justify-between text-center">
      <ButtonLoadSpinner type="submit" loading={loading}>{$t('register')}</ButtonLoadSpinner>
      <TextLink><a href="/login">{$t('back')}</a></TextLink>
    </div>
  </form>

  {#if loading}
    <Spinner />
  {/if}

  {#if showSuccess}
    <p class="bg-green-600 text-black text-center font-bold text-sm whitespace-pre-line px-4 py-2 rounded shadow-lg z-50 transition-opacity duration-300">
      {successMsg}
    </p>
  {/if}

  {#if showError}
    <div class="bg-red-600 text-white text-center font-bold text-sm whitespace-pre-line px-4 py-2 rounded shadow-lg z-50 transition-opacity duration-300">
      {errorMsg}
    </div>
  {/if}
</div>
