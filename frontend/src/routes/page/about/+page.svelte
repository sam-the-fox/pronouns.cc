<script lang="ts">
  //@ts-ignore
  import { html } from "./about.md";
</script>

<svelte:head>
  <title>About & contact - pronouns.cc</title>
  <meta name="description" content="pronouns.cc's about and contact page" />
  <meta property="og:description" content="pronouns.cc's about and contact page" />
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
