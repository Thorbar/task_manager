<script lang="ts">
  import LanguageSelector from '$lib/components/LanguageSelector.svelte';
  import CookieBanner from '$lib/components/CookieBanner.svelte';
  import '../app.css';
</script>
<!--   
    Main layout component for the app.
    - Includes a header with the language selector centered for user locale preferences.
    - Wraps page content inside <main> with padding and max-width constraints.
    - Renders a cookie consent banner to comply with privacy regulations.
    This layout is applied globally, providing consistent UI elements across pages.   
  -->
<header class="bg-slate-500 shadow-md py-4 px-8 w-320">
  <div class="max-w-7xl mx-auto flex justify-center">
    <LanguageSelector />
  </div>
</header>

<main class="max-w-7xl mx-auto px-8 py-10">
  <slot />
  <CookieBanner />
</main>
