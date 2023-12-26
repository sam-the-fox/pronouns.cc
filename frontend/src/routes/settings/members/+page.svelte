<script lang="ts">
  import type { Member } from "$lib/api/entities";
  import { Alert, ListGroup, ListGroupItem } from "@sveltestrap/sveltestrap";
  import type { PageData } from "./$types";

  export let data: PageData;

  let members: Member[];
  $: members = data.user.list_private ? data.members : data.members.filter((m) => m.unlisted);
</script>

<h1>Hidden members ({members.length})</h1>

{#if data.user.list_private}
  <Alert color="secondary" fade={false}>
    Note that your member list is hidden, so <em>all</em> of your members are shown here.
  </Alert>
{/if}

{#if members.length > 0}
  <ListGroup>
    {#each members as member}
      <ListGroupItem tag="a" href="/@{data.user.name}/{member.name}">
        <strong>{member.display_name ?? member.name}</strong>
        <span class="text-body-secondary">({member.name})</span>
      </ListGroupItem>
    {/each}
  </ListGroup>
{:else}
  <Alert color="secondary" fade={false}>You have no unlisted members!</Alert>
{/if}
