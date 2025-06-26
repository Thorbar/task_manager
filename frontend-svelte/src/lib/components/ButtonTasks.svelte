<script lang="ts">
  export let type: 'button' | 'submit' | 'reset' = 'button';
</script>

<svelte:options accessors={true} />

<div class="text-center bg-yellow-100 hover:font-bold hover:bg-yellow-200 hover:text-blue-500 transition duration-300 transform hover:scale-105">
  <button type={type} on:click {...$$restProps}>
    <slot />
  </button>
</div>
