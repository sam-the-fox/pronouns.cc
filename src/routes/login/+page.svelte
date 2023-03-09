<script lang="ts">
  import type { PageData } from "./$types";

  export let data: PageData;
</script>

<svelte:head>
  <title>Login - pronouns.cc</title>
</svelte:head>
<div>
  <a href={data.discord}>Login with Discord</a>
</div>
