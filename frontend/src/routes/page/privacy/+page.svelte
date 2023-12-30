<script lang="ts">
  // Ignoring the TS error here, because this file imports fine, typescript just chokes on markdown files
  // eslint-disable-next-line
  //@ts-ignore
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
