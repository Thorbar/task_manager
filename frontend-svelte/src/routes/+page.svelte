<script lang="ts">
  import { t } from 'svelte-i18n';
  import { goto } from '$app/navigation';
  import Button from '$lib/components/Button.svelte';

  const goLogin = () => {
    console.log('Botón clicado, navegando a /login');
    goto('/login');
  };
</script>

<!-- 
    Main landing page of the app.
    Displays a welcome message and a prominent button to navigate to the login page.
    Utilizes i18n for localized strings.
    The goLogin function handles the button click by navigating programmatically to /login.   
-->

<section class="w-full flex flex-col justify-center items-center text-center">
  <h1 class="text-4xl sm:text-5xl font-extrabold uppercase tracking-wide text-green-600 drop-shadow-lg mb-8">
    {$t('welcome')}
  </h1> 
    <Button on:click={goLogin}>{$t('login')}</Button>
</section>

