<script lang="ts">
  import { html } from "./about.md";
</script>

<svelte:head>
  <title>About & contact - pronouns.cc</title>
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
