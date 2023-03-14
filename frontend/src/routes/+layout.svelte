<script lang="ts">
  import "bootstrap/dist/css/bootstrap.min.css";
  import "bootstrap-icons/font/bootstrap-icons.css";

  import Navigation from "./nav/Navigation.svelte";
  import type { LayoutData } from "./$types";
  import { version } from "$app/environment";
  import { toastStore } from "$lib/toast";
  import Toast from "$lib/components/Toast.svelte";

  export let data: LayoutData;

  const versionParts = version.split("-");
  let commit = "[unknown]";
  if (versionParts.length >= 3) commit = versionParts[2].slice(1);

  const versionMismatch = data.git_commit !== commit && data.git_commit !== "[unknown]";
</script>

<svelte:head>
  <meta property="og:type" content="website" />
  <meta name="theme-color" content="#aa8ed6" />
  <meta property="og:site_name" content="pronouns.cc" />
</svelte:head>

<Navigation />
<div class="container">
  <slot />
  <div class="position-absolute top-0 end-0">
    {#each $toastStore as toast}
      <Toast header={toast.header} body={toast.body} />
    {/each}
  </div>
  <footer>
    <hr />
    <p>
      pronouns.cc <a href="https://codeberg.org/u1f320/pronouns.cc/commit/{commit}">{version}</a>
      {#if versionMismatch}
        (backend: <a href="https://codeberg.org/u1f320/pronouns.cc/commit/{data.git_commit}">
          {data.git_commit}
        </a>)
      {/if} &middot;
      <a href="/page/terms">Terms of service</a>
      - <a href="/page/privacy">Privacy policy</a>
    </p>
  </footer>
</div>
