<script lang="ts">
  import { marked } from "marked";
  import sanitizeHtml from "sanitize-html";

  import type { PageData } from "./$types";

  import { Alert, Icon } from "sveltestrap";
  import FieldCard from "$lib/components/FieldCard.svelte";
  import StatusIcon from "$lib/components/StatusIcon.svelte";
  import PronounLink from "$lib/components/PronounLink.svelte";
  import PartialMemberCard from "$lib/components/PartialMemberCard.svelte";
  import FallbackImage from "$lib/components/FallbackImage.svelte";
  import { userStore } from "$lib/store";
  import { pronounDisplay, userAvatars, WordStatus } from "$lib/api/entities";
  import { PUBLIC_BASE_URL } from "$env/static/public";

  export let data: PageData;

  let bio: string | null;
  $: bio = data.bio ? sanitizeHtml(marked.parse(data.bio)) : null;

  const favNames = data.names.filter((entry) => entry.status === WordStatus.Favourite);
  const favPronouns = data.pronouns.filter((entry) => entry.status === WordStatus.Favourite);
</script>

<svelte:head>
  <title>@{data.name} - pronouns.cc</title>

  <meta
    property="og:title"
    content={data.display_name ? `${data.display_name} (@${data.name})` : `@${data.name}`}
  />
  <meta property="og:url" content="{PUBLIC_BASE_URL}/@{data.name}" />

  {#if data.avatar}
    <meta property="og:image" content={userAvatars(data)[0]} />
  {/if}

  {#if favNames.length !== 0 && favPronouns.length !== 0}
    <meta
      property="og:description"
      content="@{data.name} goes by {favNames.map((x) => x.value).join(', ')} and uses {favPronouns
        .map((x) => pronounDisplay(x))
        .join(', ')} pronouns."
    />
  {:else if favNames.length !== 0}
    <meta
      property="og:description"
      content="@{data.name} goes by {favNames.map((x) => x.value).join(', ')}."
    />
  {:else if favPronouns.length !== 0}
    <meta
      property="og:description"
      content="@{data.name} uses {favPronouns.map((x) => pronounDisplay(x)).join(', ')} pronouns."
    />
  {:else if data.bio && data.bio !== ""}
    <meta
      property="og:description"
      content="{data.bio.slice(0, 500)}${data.bio.length > 500 ? '…' : ''}"
    />
  {:else}
    <meta property="og:description" content="@{data.name} on pronouns.cc" />
  {/if}
</svelte:head>

<div class="container">
  {#if $userStore && $userStore.id === data.id}
    <Alert color="secondary" fade={false}>
      You are currently viewing your <strong>public</strong> profile.
      <br /><a href="/edit/profile">Edit your profile</a>
    </Alert>
  {/if}
  <div class="grid row-gap-3">
    <div class="row">
      <div class="col-md text-center">
        <FallbackImage width={200} urls={userAvatars(data)} alt="Avatar for @{data.name}" />
      </div>
      <div class="col-md">
        {#if data.display_name}
          <div>
            <h2>{data.display_name}</h2>
            <h5 class="text-body-secondary">@{data.name}</h5>
          </div>
        {:else}
          <h2>@{data.name}</h2>
        {/if}
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
          <h4>Names</h4>
          <ul class="list-unstyled">
            {#each data.names as name}
              <li><StatusIcon status={name.status} /> {name.value}</li>
            {/each}
          </ul>
        </div>
      {/if}
      {#if data.pronouns.length > 0}
        <div class="col-md">
          <h4>Pronouns</h4>
          <ul class="list-unstyled">
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
    {#if data.members.length > 0}
      <div class="row">
        <div class="col">
          <hr />
          <h2>Members</h2>
        </div>
      </div>
      <div class="row row-cols-1 row-cols-sm-2 row-cols-md-4 text-center">
        {#each data.members as member}
          <PartialMemberCard user={data} {member} />
        {/each}
      </div>
    {/if}
  </div>
</div>
