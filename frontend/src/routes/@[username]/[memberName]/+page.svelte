<script lang="ts">
  import { marked } from "marked";
  import sanitizeHtml from "sanitize-html";

  import FieldCard from "$lib/components/FieldCard.svelte";

  import type { PageData } from "./$types";
  import StatusIcon from "$lib/components/StatusIcon.svelte";
  import PronounLink from "$lib/components/PronounLink.svelte";
  import FallbackImage from "$lib/components/FallbackImage.svelte";
  import { Button, Icon } from "sveltestrap";

  export let data: PageData;

  let bio: string | null;
  $: bio = data.bio ? sanitizeHtml(marked.parse(data.bio)) : null;
</script>

<svelte:head>
  <title>{data.display_name ?? data.name} - @{data.user.name} - pronouns.cc</title>
</svelte:head>

<div class="container">
  <div>
    <Button color="secondary" href="/@{data.user.name}">
      <Icon name="arrow-left" /> Back to {data.user.display_name ?? data.user.name}
    </Button>
  </div>
  <div class="grid">
    <div class="row">
      <div class="col-md text-center">
        <FallbackImage urls={data.avatar_urls} alt="Avatar for @{data.name}" />
      </div>
      <div class="col-md">
        <h2>{data.display_name ?? data.name}</h2>
        <h5 class="text-body-secondary">{data.name} (@{data.user.name})</h5>
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
                <PronounLink {pronouns} />
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
