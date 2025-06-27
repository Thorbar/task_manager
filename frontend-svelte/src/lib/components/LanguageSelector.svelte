<script lang="ts">
  import { t } from 'svelte-i18n';
  import { locale } from 'svelte-i18n';
  import { goto } from '$app/navigation';

  const availableLocales = [
    { code: 'es', label: 'Español', flag: '/flags/es.jpg' },
    { code: 'ca', label: 'Català', flag: '/flags/ca.jpg' },
    { code: 'en', label: 'English', flag: '/flags/en.jpg' },
    { code: 'ja', label: '日本語', flag: '/flags/ja.jpg' }
  ];

  let isOpen = false;

  function selectLocale(code: string) {
    locale.set(code);
    isOpen = false;
  }

  const goHome = () => goto('/');
</script>

<div class="flex items-center space-x-2  ">
  <button
    on:click={goHome}
    aria-label="Ir al menú principal"
    class="text-2xl hover:text-green-600 transition-colors"
    title="Ir al menú principal"
  >
    🌐
  </button>

 <label for="lang" class="w-16 font-semibold text-gray-700 select-none">{$t('language')}:</label>
 
  <!-- Dropdown personalizado -->
  <div class="relative">
    <button
      class="flex items-center space-x-2   rounded-md p-1 px-3 focus:outline-none focus:ring-2 focus:ring-white"
      aria-haspopup="listbox"
      aria-expanded={isOpen}
      on:click={() => (isOpen = !isOpen)}
    >
      {#each availableLocales as loc}
        {#if loc.code === $locale}
          <img src={loc.flag} alt={loc.label} class="w-5 h-5 rounded-sm" />
          <span class="ml-1">{loc.label}</span>
        {/if}
      {/each}
      <svg
        class="w-4 h-4 ml-2"
        xmlns="http://www.w3.org/2000/svg"
        fill="none"
        viewBox="0 0 24 24"
        stroke="currentColor"
      >
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
      </svg>
    </button>

    {#if isOpen}
      <ul 
        role="listbox"
        class="absolute z-10 mt-1 w-full bg-white border border-gray-300 rounded-md shadow-lg max-h-48 overflow-auto"
      >
      {#each availableLocales as loc}
        <li
          role="option"
          tabindex="0"
          aria-selected={loc.code === $locale}
          class="cursor-pointer flex items-center space-x-2 p-2 hover:bg-white-100"
          on:click={() => selectLocale(loc.code)}
          on:keydown={(e) => {
            if (e.key === 'Enter' || e.key === ' ') {
              e.preventDefault();
              selectLocale(loc.code);
            }
          }}
        >
          <img src={loc.flag} alt={loc.label} class="w-5 h-5 rounded-sm" />
          <span>{loc.label}</span>
        </li>
      {/each}
      </ul>
    {/if}
  </div>
</div>
