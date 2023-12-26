<script lang="ts">
  import "bootstrap/dist/css/bootstrap.min.css";
  import "@fontsource/firago/400.css";
  import "@fontsource/firago/400-italic.css";
  import "@fontsource/firago/700.css";
  import "@fontsource/firago/700-italic.css";
  import "bootstrap-icons/font/bootstrap-icons.css";
  import "./main.css";

  import Navigation from "./nav/Navigation.svelte";
  import type { LayoutData } from "./$types";
  import { version } from "$app/environment";
  import { settingsStore } from "$lib/store";
  import { toastStore } from "$lib/toast";
  import Toast from "$lib/components/Toast.svelte";
  import { Alert, Icon } from "@sveltestrap/sveltestrap";
  import { apiFetchClient } from "$lib/api/fetch";
  import type { Settings } from "$lib/api/entities";
  import { renderUnsafeMarkdown } from "$lib/utils";

  export let data: LayoutData;

  const versionParts = version.split("-");
  let commit = "[unknown]";
  if (versionParts.length >= 3) commit = versionParts[2].slice(1);

  const versionMismatch = data.git_commit !== commit && data.git_commit !== "[unknown]";

  const readNotice = async () => {
    try {
      const resp = await apiFetchClient<Settings>(
        "/users/@me/settings",
        "PATCH",
        { read_global_notice: data.notice!.id },
        2,
      );
      settingsStore.set({ current: true, settings: resp });
    } catch (e) {
      console.log("updating settings:", e);
    }
  };
</script>

<svelte:head>
  <meta property="og:type" content="website" />
  <meta name="theme-color" content="#aa8ed6" />
  <meta property="og:site_name" content="pronouns.cc" />
  <script defer data-domain="pronouns.cc" src="https://plausible.pronouns.cc/js/script.js"></script>
</svelte:head>

<div class="d-flex flex-column min-vh-100">
  <div class="flex-grow-1">
    <Navigation commit={data.git_commit} />
    <div class="container">
      {#if data.notice && $settingsStore.current && data.notice.id > $settingsStore.settings.read_global_notice}
        <Alert color="secondary" isOpen={true} toggle={() => readNotice()}>
          {@html renderUnsafeMarkdown(data.notice.notice)}
        </Alert>
      {/if}

      <slot />
      <div class="position-absolute top-0 start-50 translate-middle-x">
        {#each $toastStore as toast}
          <Toast header={toast.header} body={toast.body} />
        {/each}
      </div>
    </div>
  </div>

  <footer class="container">
    <hr />
    <div class="d-flex flex-column flex-md-row">
      <p>
        pronouns.cc <a href="https://codeberg.org/pronounscc/pronouns.cc/commit/{commit}"
          >{version}</a
        >
        {#if versionMismatch}
          (backend: <a href="https://codeberg.org/pronounscc/pronouns.cc/commit/{data.git_commit}"
            >{data.git_commit}</a
          >)
        {/if} &middot;
        <a href="https://status.pronouns.cc/" target="_blank">Status</a> &middot;
        <a href="/page/about">About & contact</a> &middot;
        <a href="/page/changelog">Changelog</a> &middot;
        <Icon name="cash" aria-hidden />
        <a href="https://liberapay.com/u1f320/" target="_blank">Donate</a>
        &middot;
        <a href="/page/terms">Terms of service</a> &middot;
        <a href="/page/privacy">Privacy policy</a>
      </p>
      <p class="ms-auto">
        Users: <strong>{data.users.total}</strong> &middot; Members: <strong>{data.members}</strong>
      </p>
    </div>
  </footer>
</div>
