<script lang="ts">
  import { getContext } from "svelte";
  import type { Writable } from "svelte/store";
  import { Alert, ButtonGroup, Input } from "sveltestrap";
  import type { PageData } from "./$types";

  import type { Member, PrideFlag } from "$lib/api/entities";
  import IconButton from "$lib/components/IconButton.svelte";
  import FlagButton from "$lib/components/edit/FlagButton.svelte";

  export let data: PageData;

  const member = getContext<Writable<Member>>("member");

  let flagSearch = "";
  let filteredFlags: PrideFlag[];
  $: filteredFlags = filterFlags(flagSearch, data.flags);

  const filterFlags = (search: string, flags: PrideFlag[]) => {
    return (
      search
        ? flags.filter((flag) => flag.name.toLocaleLowerCase().includes(search.toLocaleLowerCase()))
        : flags
    ).slice(0, 25);
  };

  const addFlag = (flag: PrideFlag) => {
    $member.flags = [...$member.flags, flag];
  };

  const moveFlag = (index: number, up: boolean) => {
    if (up && index == 0) return;
    if (!up && index == $member.flags.length - 1) return;

    const newIndex = up ? index - 1 : index + 1;

    const temp = $member.flags[index];
    $member.flags[index] = $member.flags[newIndex];
    $member.flags[newIndex] = temp;
    $member.flags = [...$member.flags];
  };

  const removeFlag = (index: number) => {
    $member.flags.splice(index, 1);
    $member.flags = [...$member.flags];
  };
</script>

<div>
  {#each $member.flags as _, index}
    <ButtonGroup class="m-1">
      <IconButton
        icon="chevron-left"
        color="secondary"
        tooltip="Move flag to the left"
        click={() => moveFlag(index, true)}
      />
      <IconButton
        icon="chevron-right"
        color="secondary"
        tooltip="Move flag to the right"
        click={() => moveFlag(index, false)}
      />
      <FlagButton
        flag={$member.flags[index]}
        tooltip="Remove this flag from your profile"
        on:click={() => removeFlag(index)}
      />
    </ButtonGroup>
  {/each}
</div>
<hr />
<div class="row">
  <div class="col-md">
    <Input placeholder="Filter flags" bind:value={flagSearch} disabled={data.flags.length === 0} />
    <div class="p-2">
      {#each filteredFlags as flag (flag.id)}
        <FlagButton {flag} tooltip="Add this flag to your profile" on:click={() => addFlag(flag)} />
      {:else}
        {#if data.flags.length === 0}
          You haven't uploaded any flags yet.
        {:else}
          There are no flags matching your search <strong>{flagSearch}</strong>.
        {/if}
      {/each}
    </div>
  </div>
  <div class="col-md">
    <Alert color="secondary" fade={false}>
      {#if data.flags.length === 0}
        <p><strong>Why can't I see any flags?</strong></p>
        <p>
          There are thousands of pride flags, and it would be impossible to bundle all of them by
          default. Many labels also have multiple different flags that are favoured by different
          people. Because of this, there are no flags available by default--instead, you can upload
          flags in your <a href="/settings/flags">settings</a>. Your main profile and your member
          profiles can all have different flags.
        </p>
      {:else}
        To upload and delete flags, go to your <a href="/settings/flags">settings</a>.
      {/if}
    </Alert>
  </div>
</div>
