<script lang="ts">
  export let href: string;
  export let plain: boolean = false;
</script>

{#if plain}
  <a {href} class="hover:text-sky-500 dark:hover:text-sky-400"><slot /></a>
{:else}
  <li>
    <a {href} class="hover:text-sky-500 dark:hover:text-sky-400"><slot /></a>
  </li>
{/if}
