<script lang="ts">
  export let href = '';
</script>

<div
  class={`text-sm font-bold text-gray-600 underline hover:text-blue-500 transition duration-200 transform hover:scale-110 ${href}`}
>
  <slot />
</div>
