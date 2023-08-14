<script lang="ts">
  import { getContext } from "svelte";
  import type { Writable } from "svelte/store";

  import type { Member } from "$lib/api/entities";
  import IconButton from "$lib/components/IconButton.svelte";

  const member = getContext<Writable<Member>>("member");

  let newLink = "";

  const addLink = (event: Event) => {
    event.preventDefault();

    $member.links = [...$member.links, newLink];
    newLink = "";
  };

  const removeLink = (index: number) => {
    $member.links.splice(index, 1);
    $member.links = [...$member.links];
  };

  const moveLink = (index: number, up: boolean) => {
    if (up && index == 0) return;
    if (!up && index == $member.links.length - 1) return;

    const newIndex = up ? index - 1 : index + 1;

    const temp = $member.links[index];
    $member.links[index] = $member.links[newIndex];
    $member.links[newIndex] = temp;
    $member.links = [...$member.links];
  };
</script>

{#each $member.links as _, index}
  <div class="input-group m-1">
    <IconButton
      icon="chevron-up"
      color="secondary"
      tooltip="Move link up"
      click={() => moveLink(index, true)}
    />
    <IconButton
      icon="chevron-down"
      color="secondary"
      tooltip="Move link down"
      click={() => moveLink(index, false)}
    />
    <input type="text" class="form-control" bind:value={$member.links[index]} />
    <IconButton
      color="danger"
      icon="trash3"
      tooltip="Remove link"
      click={() => removeLink(index)}
    />
  </div>
{/each}
<form class="input-group m-1" on:submit={addLink}>
  <input type="text" class="form-control" bind:value={newLink} />
  <IconButton type="submit" color="success" icon="plus" tooltip="Add link" />
</form>
