<script lang="ts">
  import FieldCard from "$lib/components/FieldCard.svelte";

  import type { PageData } from "./$types";
  import PronounLink from "$lib/components/PronounLink.svelte";
  import FallbackImage from "$lib/components/FallbackImage.svelte";
  import { Alert, Button, Icon } from "sveltestrap";
  import {
    memberAvatars,
    pronounDisplay,
    type CustomPreferences,
    type FieldEntry,
    type Pronoun,
  } from "$lib/api/entities";
  import { PUBLIC_BASE_URL } from "$env/static/public";
  import { userStore } from "$lib/store";
  import { renderMarkdown } from "$lib/utils";
  import ReportButton from "../ReportButton.svelte";
  import ProfileLink from "../ProfileLink.svelte";
  import StatusLine from "$lib/components/StatusLine.svelte";
  import defaultPreferences from "$lib/api/default_preferences";

  export let data: PageData;

  let bio: string | null;
  $: bio = renderMarkdown(data.bio);

  let mergedPreferences: CustomPreferences;
  $: mergedPreferences = Object.assign(defaultPreferences, data.user.custom_preferences);

  let favNames: FieldEntry[];
  $: favNames = data.names.filter(
    (entry) => (mergedPreferences[entry.status] || { favourite: false }).favourite,
  );
  let favPronouns: Pronoun[];
  $: favPronouns = data.pronouns.filter(
    (entry) => (mergedPreferences[entry.status] || { favourite: false }).favourite,
  );

  let profileEmpty = false;
  $: profileEmpty =
    data.names.length === 0 &&
    data.pronouns.length === 0 &&
    data.fields.length === 0 &&
    (!data.bio || data.bio.length === 0);
</script>

<div class="container">
  {#if $userStore && $userStore.id === data.user.id}
    <Alert color="secondary" fade={false}>
      You are currently viewing the <strong>public</strong> profile of {data.display_name ??
        data.name}.
      <br /><a href="/edit/member/{data.id}">Edit profile</a>
    </Alert>
  {/if}
  <div class="m-3">
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
        <p class="fs-5 text-body-secondary">{data.name} (@{data.user.name})</p>
        {#if profileEmpty && $userStore?.id === data.user.id}
          <hr />
          <p>
            <em>
              This member's profile is empty! You can customize it by going to the <a
                href="/edit/member/{data.id}">edit member</a
              > page.</em
            > <span class="text-muted">(only you can see this)</span>
          </p>
        {:else if bio}
          <hr />
          <p>{@html bio}</p>
        {/if}
      </div>
      {#if data.links.length > 0}
        <div class="col-md d-flex align-items-center">
          <ul class="list-unstyled">
            {#each data.links as link}
              <ProfileLink {link} />
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
              <li>
                <StatusLine preferences={data.user.custom_preferences} status={name.status}
                  >{name.value}</StatusLine
                >
              </li>
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
                <StatusLine preferences={data.user.custom_preferences} status={pronouns.status}
                  ><PronounLink {pronouns} /></StatusLine
                >
              </li>
            {/each}
          </ul>
        </div>
      {/if}
      {#each data.fields as field}
        <div class="col">
          <FieldCard preferences={data.user.custom_preferences} {field} />
        </div>
      {/each}
    </div>
    {#if $userStore && $userStore.id !== data.user.id}
      <div class="row">
        <ReportButton subject="member" reportUrl="/members/{data.id}/reports" />
      </div>
    {/if}
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
  <meta name="description" content="{data.name} (@{data.user.name}) on pronouns.cc" />
  <meta name="robots" content="noindex" />

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
