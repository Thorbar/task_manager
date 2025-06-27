<script lang="ts">
  export let type: 'button' | 'submit' | 'reset' = 'button';
</script>

<svelte:options accessors={true} />

<button
  type={type}
  class="bg-amber-700 hover:bg-cyan-600 hover:text-black text-white font-bold py-2 px-6 rounded-lg shadow-lg transition duration-300 transform hover:scale-105"
  on:click
  {...$$restProps}
>
  <slot />
</button>
