<script lang="ts">
  import { onMount } from 'svelte';
  import { fade } from 'svelte/transition';
  import { t } from 'svelte-i18n';

  let showBanner = false;

  onMount(() => {
    const accepted = localStorage.getItem('cookiesAccepted');
    if (!accepted) {
      showBanner = true;
    }
  });

  function acceptCookies() {
    localStorage.setItem('cookiesAccepted', 'true');
    showBanner = false;
  }
</script> 

{#if showBanner} 
  <div class="fixed inset-0 z-50 flex items-center justify-center bg-black/40" transition:fade>
    <div class="bg-white h-40 max-w-lg w-[90%] h-30 rounded-2xl shadow-xl p-6 text-gray-800 text-sm space-y-4 text-center justify-center flex flex-col items-center">
        <div class="text-sm text-gray-800 h-23">
            {#each $t('cookiesAlert').split('\n') as line}
                <p>{line}</p>
            {/each}
        </div>
        
        <div class="bg-green-600 text-white w-[20%] items-center text-center px-4 py-2 rounded-xl hover:bg-green-700 transition">
            <button on:click={acceptCookies} >     
                {$t('accept')}
            </button>
        </div>        
    </div>
  </div>
{/if}
