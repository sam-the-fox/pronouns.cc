<script lang="ts">
  import type { PageData } from "./$types";

  export let data: PageData;
</script>

<svelte:head>
  <title>Log in with Discord - pronouns.cc</title>
</svelte:head>

<h1>Log in with Discord</h1>

{#if data.error}

{:else}

{/if}