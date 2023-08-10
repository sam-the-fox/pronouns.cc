<script lang="ts">
  import { getContext } from "svelte";
  import type { Writable } from "svelte/store";

  import type { MeUser } from "$lib/api/entities";
  import IconButton from "$lib/components/IconButton.svelte";

  const user = getContext<Writable<MeUser>>("user");

  let newLink = "";

  const addLink = (event: Event) => {
    event.preventDefault();

    $user.links = [...$user.links, newLink];
    newLink = "";
  };

  const removeLink = (index: number) => {
    $user.links.splice(index, 1);
    $user.links = [...$user.links];
  };

  const moveLink = (index: number, up: boolean) => {
    if (up && index == 0) return;
    if (!up && index == $user.links.length - 1) return;

    const newIndex = up ? index - 1 : index + 1;

    const temp = $user.links[index];
    $user.links[index] = $user.links[newIndex];
    $user.links[newIndex] = temp;
    $user.links = [...$user.links];
  };
</script>

{#each $user.links as _, index}
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
    <input type="text" class="form-control" bind:value={$user.links[index]} />
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
