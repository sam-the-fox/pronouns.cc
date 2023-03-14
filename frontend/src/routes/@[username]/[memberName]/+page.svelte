<script lang="ts">
  import { marked } from "marked";
  import sanitizeHtml from "sanitize-html";

  import FieldCard from "$lib/components/FieldCard.svelte";

  import type { PageData } from "./$types";
  import StatusIcon from "$lib/components/StatusIcon.svelte";
  import PronounLink from "$lib/components/PronounLink.svelte";
  import FallbackImage from "$lib/components/FallbackImage.svelte";
  import { Alert, Button, Icon } from "sveltestrap";
  import { memberAvatars, pronounDisplay, WordStatus } from "$lib/api/entities";
  import { PUBLIC_BASE_URL } from "$env/static/public";
  import { userStore } from "$lib/store";

  export let data: PageData;

  let bio: string | null;
  $: bio = data.bio ? sanitizeHtml(marked.parse(data.bio)) : null;

  const favNames = data.names.filter((entry) => entry.status === WordStatus.Favourite);
  const favPronouns = data.pronouns.filter((entry) => entry.status === WordStatus.Favourite);
</script>

<div class="container">
  {#if $userStore && $userStore.id === data.user.id}
    <Alert color="secondary" fade={false}>
      You are currently viewing the <strong>public</strong> profile of {data.display_name ??
        data.name}.
      <br /><a href="/edit/member/{data.id}">Edit profile</a>
    </Alert>
  {/if}
  <div>
    <Button color="secondary" href="/@{data.user.name}">
      <Icon name="arrow-left" /> Back to {data.user.display_name ?? data.user.name}
    </Button>
  </div>
  <div class="grid">
    <div class="row">
      <div class="col-md-4 text-center">
        <FallbackImage width={200} urls={memberAvatars(data)} alt="Avatar for @{data.name}" />
      </div>
      <div class="col-md">
        <h2>{data.display_name ?? data.name}</h2>
        <h5 class="text-body-secondary">{data.name} (@{data.user.name})</h5>
        {#if bio}
          <hr />
          <p>{@html bio}</p>
        {/if}
      </div>
      {#if data.links.length > 0}
        <div class="col-md d-flex align-items-center">
          <ul class="list-unstyled">
            {#each data.links as link}
              <li><Icon name="globe" /> <a href={link}>{link}</a></li>
            {/each}
          </ul>
        </div>
      {/if}
    </div>
    <div class="row row-cols-1 row-cols-sm-2 row-cols-md-3">
      {#if data.names.length > 0}
        <div class="col-md">
          <h3>Names</h3>
          <ul class="list-unstyled fs-5">
            {#each data.names as name}
              <li><StatusIcon status={name.status} /> {name.value}</li>
            {/each}
          </ul>
        </div>
      {/if}
      {#if data.pronouns.length > 0}
        <div class="col-md">
          <h3>Pronouns</h3>
          <ul class="list-unstyled fs-5">
            {#each data.pronouns as pronouns}
              <li>
                <StatusIcon status={pronouns.status} />
                <PronounLink {pronouns} />
              </li>
            {/each}
          </ul>
        </div>
      {/if}
      {#each data.fields as field}
        <div class="col">
          <FieldCard {field} />
        </div>
      {/each}
    </div>
  </div>
</div>

<svelte:head>
  <title>{data.display_name ?? data.name} - @{data.user.name} - pronouns.cc</title>

  <meta
    property="og:title"
    content={data.display_name
      ? `${data.display_name} (${data.name})`
      : `${data.name} (@${data.user.name})`}
  />
  <meta property="og:url" content="{PUBLIC_BASE_URL}/@{data.user.name}/{data.name}" />

  {#if data.avatar}
    <meta property="og:image" content={memberAvatars(data)[0]} />
  {/if}

  {#if favNames.length !== 0 && favPronouns.length !== 0}
    <meta
      property="og:description"
      content="{data.display_name ?? data.name} goes by {favNames
        .map((x) => x.value)
        .join(', ')} and uses {favPronouns.map((x) => pronounDisplay(x)).join(', ')} pronouns."
    />
  {:else if favNames.length !== 0}
    <meta
      property="og:description"
      content="{data.display_name ?? data.name} goes by {favNames.map((x) => x.value).join(', ')}."
    />
  {:else if favPronouns.length !== 0}
    <meta
      property="og:description"
      content="{data.display_name ?? data.name} uses {favPronouns
        .map((x) => pronounDisplay(x))
        .join(', ')} pronouns."
    />
  {:else if data.bio && data.bio !== ""}
    <meta
      property="og:description"
      content="{data.bio.slice(0, 500)}{data.bio.length > 500 ? '…' : ''}"
    />
  {:else}
    <meta property="og:description" content="{data.name} (@{data.user.name}) on pronouns.cc" />
  {/if}
</svelte:head>
