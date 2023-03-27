<script lang="ts">
  import { html } from "./privacy.md";
</script>

<svelte:head>
  <title>Privacy policy - pronouns.cc</title>
  <meta name="description" content="pronouns.cc's privacy policy" />
  <meta property="og:description" content="pronouns.cc's privacy policy" />
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
