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
        const errorData = await res.json();

        if (res.status === 409 && errorData?.error?.includes('exists')) {
          showErrorMessage($t('userAlreadyExists'));
        } else if (res.status === 400 && errorData?.error?.includes('password')) {
          showErrorMessage($t('wrongPassword'));
        } else {
          console.log('4');
          showErrorMessage($t('registeroFailed'));
        }

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

<section class="flex flex-col items-center">
    <h2 class="text-3xl font-extrabold uppercase tracking-wide text-green-600 drop-shadow-lg mb-10">{$t('register')}</h2>

  <form on:submit|preventDefault={registerUser} class="space-y-3 w-full flex flex-col items-center">
    <div class="font-bold">
      <label for="name" >{$t('name')}</label>
      <InputWrapper>
        <input id="name" bind:value={name} type="text" required class="w-full" />
      </InputWrapper>
    </div>

    <div class="font-bold">
      <label for="surname" >{$t('surname')}</label>
      <InputWrapper>
        <input id="surname" bind:value={surname} type="text" required class="w-full" />
      </InputWrapper>
    </div>

    <div class="font-bold">
      <label for="email" >{$t('email')}</label>
      <InputWrapper>
        <input id="email" bind:value={email} type="email" required class="w-full" />
      </InputWrapper>
    </div>

    <div class="w-55 font-bold">
      <label for="username">{$t('username')}</label>
      <InputWrapper>
        <input id="username" bind:value={username} type="text" required class="w-full" />
      </InputWrapper>
    </div>

    <div class="w-55 font-bold">
      <label for="password">{$t('password')}</label>
      <InputWrapper>
        <input id="password" bind:value={password} type="password" required class="w-full" />
      </InputWrapper>
    </div>

    <div class="h-15 w-55 font-bold">
      <label for="confirmPassword">{$t('confirmPassword')}</label>
      <InputWrapper>
        <input id="confirmPassword" bind:value={confirmPassword} type="password" required class="w-full" />
      </InputWrapper>
    </div>

    <div class="flex justify-between items-center space-x-4">
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
</section>
