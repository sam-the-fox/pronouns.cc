<script lang="ts">
  import { marked } from "marked";
  import sanitizeHtml from "sanitize-html";

  import FieldCard from "$lib/components/FieldCard.svelte";

  import type { PageData } from "./$types";
  import StatusIcon from "$lib/components/StatusIcon.svelte";
    import PronounLink from "$lib/components/PronounLink.svelte";

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
      {#if data.avatar_urls}
        <div class="col-md" />
      {/if}
      <div class="col-md">
        {#if data.display_name}
          <h2>{data.display_name}</h2>
          <h4>@{data.name}</h4>
        {:else}
          <h2>@{data.name}</h2>
        {/if}
        <hr />
        {#if bio}
          <p>{@html bio}</p>
        {/if}
      </div>
      {#if data.links}
        <div class="col-md">
          <ul>
            {#each data.links as link}
              <li><a href={link}>{link}</a></li>
            {/each}
          </ul>
        </div>
      {/if}
    </div>
    <div class="row">
      {#if data.names}
        <div class="col-md">
          <h4>Names</h4>
          <ul class="list-unstyled">
            {#each data.names as name}
              <li><StatusIcon status={name.status} /> {name.value}</li>
            {/each}
          </ul>
        </div>
      {/if}
      {#if data.pronouns}
        <div class="col-md">
          <h4>Pronouns</h4>
          <ul class="list-unstyled">
            {#each data.pronouns as pronouns}
              <li>
                <StatusIcon status={pronouns.status} />
                <PronounLink pronouns={pronouns} />
              </li>
            {/each}
          </ul>
        </div>
      {/if}
    </div>
  </div>
  {#if data.fields}
    <div class="row row-cols-1 row-cols-sm-2 row-cols-md-3">
      {#each data.fields as field}
        <div class="col">
          <FieldCard {field} />
        </div>
      {/each}
    </div>
  {/if}
</div>
