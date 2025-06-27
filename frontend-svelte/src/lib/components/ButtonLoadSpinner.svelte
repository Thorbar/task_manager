<script lang="ts">
  export let type: 'button' | 'submit' | 'reset' = 'button';
  export let loading: boolean = false;
</script>

<svelte:options accessors={true} />

<div class="bg-amber-700 hover:bg-cyan-600 hover:text-black text-white font-bold py-4 px-10 
rounded-lg shadow-lg transition duration-300 transform hover:scale-105">
  <button type={type}
    disabled={loading}>
    <slot />
  </button>
</div>

