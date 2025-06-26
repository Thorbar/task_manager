<script lang="ts">
  export let className = '';
</script>

<div class={`bg-white border rounded w-55 focus-within:ring-2 focus-within:ring-green-400 p-2 ${className}`}>
  <slot />
</div>

