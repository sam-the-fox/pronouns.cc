<script lang="ts">
  //@ts-ignore
  import { html } from "./terms.md";
</script>

<svelte:head>
  <title>Terms of service - pronouns.cc</title>
  <meta name="description" content="pronouns.cc's terms of service" />
  <meta property="og:description" content="pronouns.cc's terms of service" />
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
