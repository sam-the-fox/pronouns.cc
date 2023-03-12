<script lang="ts">
  import { marked } from "marked";
  import sanitizeHtml from "sanitize-html";

  import type { PageData } from "./$types";

  import { Icon } from "sveltestrap";
  import FieldCard from "$lib/components/FieldCard.svelte";
  import StatusIcon from "$lib/components/StatusIcon.svelte";
  import PronounLink from "$lib/components/PronounLink.svelte";
  import PartialMemberCard from "$lib/components/PartialMemberCard.svelte";
  import FallbackImage from "$lib/components/FallbackImage.svelte";

  export let data: PageData;

  let bio: string | null;
  $: bio = data.bio ? sanitizeHtml(marked.parse(data.bio)) : null;
</script>

<svelte:head>
  <title>@{data.name} - pronouns.cc</title>
</svelte:head>

<div class="container">
  <div class="grid">
    <div class="row">
      <div class="col-md text-center">
        <FallbackImage urls={data.avatar_urls} alt="Avatar for @{data.name}" />
      </div>
      <div class="col-md-8">
        {#if data.display_name}
          <h2>{data.display_name}</h2>
          <h5 class="text-body-secondary">@{data.name}</h5>
        {:else}
          <h2>@{data.name}</h2>
        {/if}
        <hr />
        {#if bio}
          <p>{@html bio}</p>
        {/if}
        {#if data.links.length > 0}
          <hr />
          <ul class="list-unstyled">
            {#each data.links as link}
              <li><Icon name="globe" /> <a href={link}>{link}</a></li>
            {/each}
          </ul>
        {/if}
      </div>
    </div>
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
