<script>
  import { t } from 'svelte-i18n';
</script>

<div class="fixed inset-0  flex items-center justify-center h-90  w-290">
  <p class="text-white bg-black rounded-lg shadow-lg text-xl font-semibold font-mono loading-dots">
    {$t('loading')}<span class="dots"></span>
  </p>
</div>
