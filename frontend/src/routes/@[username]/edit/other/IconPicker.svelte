<script lang="ts">
  import {
    ButtonDropdown,
    DropdownItem,
    DropdownMenu,
    DropdownToggle,
    Icon,
    Input,
    Tooltip,
  } from "sveltestrap";
  import icons from "../../../../icons";
  import IconButton from "$lib/components/IconButton.svelte";

  export let icon = "";

  const MAX_VISIBLE_ITEMS = 21;

  let searchBox = "";
  let showAll = false;
  let filteredIcons: string[] = [];
  $: filteredIcons = icons
    .filter((icon) => icon.toLowerCase().includes(searchBox.toLowerCase()))
    .sort()
    .slice(0, showAll ? undefined : MAX_VISIBLE_ITEMS);

  let iconCount = 0;
  $: iconCount = icons.filter((icon) =>
    icon.toLowerCase().includes(searchBox.toLowerCase()),
  ).length;

  let iconButton: HTMLElement;
</script>

<ButtonDropdown autoClose={false}>
  <Tooltip target={iconButton} placement="top">Change icon</Tooltip>
  <DropdownToggle color="secondary" caret bind:inner={iconButton}>
    <Icon name={icon} alt="Current icon" />
  </DropdownToggle>
  <DropdownMenu>
    <p class="px-2">
      <Input type="text" placeholder="Search for icons" bind:value={searchBox} />
    </p>
    <ul class="list-unstyled icon-list text-center">
      {#each filteredIcons as iconOption}
        <li class="list-inline-item">
          <IconButton
            icon={iconOption}
            tooltip={iconOption}
            outline
            border={false}
            color="secondary"
            click={() => (icon = iconOption)}
          />
        </li>
      {/each}
    </ul>
    {#if iconCount > MAX_VISIBLE_ITEMS || showAll}
      <DropdownItem divider />
      <DropdownItem toggle={true} on:click={() => (showAll = !showAll)}
        >{showAll ? "Stop showing all results" : `Show all results (${iconCount})`}</DropdownItem
      >
    {/if}
  </DropdownMenu>
</ButtonDropdown>

<style>
  .icon-list {
    max-height: 200px;
    overflow-y: auto;
  }
</style>
