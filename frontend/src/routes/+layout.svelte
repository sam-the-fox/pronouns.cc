<script lang="ts">
  import "bootstrap/dist/css/bootstrap.min.css";
  import "bootstrap-icons/font/bootstrap-icons.css";

  import Navigation from "./nav/Navigation.svelte";
  import type { LayoutData } from "./$types";
  import { version } from "$app/environment";

  export let data: LayoutData;

  const versionMismatch = data.git_commit !== version && data.git_commit !== "[unknown]";
</script>

<Navigation />
<div class="container">
  <slot />
  <footer>
    <hr />
    <p>
      pronouns.cc <a href={data.git_repository}>{version}</a>
      {#if versionMismatch}(backend: {data.git_commit}){/if} &middot;
      <a href="/page/terms">Terms of service</a>
      - <a href="/page/privacy">Privacy policy</a>
    </p>
  </footer>
</div>
