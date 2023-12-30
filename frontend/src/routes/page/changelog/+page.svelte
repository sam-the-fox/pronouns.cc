<script lang="ts">
  import { onMount } from "svelte";
  // Ignoring the TS error here, because this file imports fine, typescript just chokes on markdown files
  // eslint-disable-next-line
  //@ts-ignore
  import { html } from "./changelog.md";
  import { CURRENT_CHANGELOG } from "$lib/store";

  onMount(() => {
    localStorage.setItem("changelog-read", CURRENT_CHANGELOG);
  });
</script>

<svelte:head>
  <title>Changelog - pronouns.cc</title>
  <meta name="description" content="pronouns.cc's changelog" />
  <meta property="og:description" content="pronouns.cc's changelog" />
</svelte:head>

<div class="grid">
  <div class="row">
    <div class="col-lg-2" />
    <div class="col-lg-8">
      {@html html}
    </div>
    <div class="col-lg-2" />
  </div>
</div>
