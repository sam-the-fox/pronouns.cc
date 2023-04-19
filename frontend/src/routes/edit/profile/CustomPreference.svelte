<script lang="ts">
  import { PreferenceSize, type CustomPreference } from "$lib/api/entities";
  import IconButton from "$lib/components/IconButton.svelte";
  import {
    ButtonDropdown,
    DropdownItem,
    DropdownMenu,
    DropdownToggle,
    Icon,
    Input,
    InputGroup,
    Tooltip,
  } from "sveltestrap";
  import icons from "../../../icons";

  export let preference: CustomPreference;
  export let remove: () => void;

  let iconButton: HTMLElement;
  let sizeButton: HTMLElement;

  const toggleMuted = () => (preference.muted = !preference.muted);
  const toggleFavourite = () => (preference.favourite = !preference.favourite);

  let searchBox = "";
  let filteredIcons: string[] = [];
  $: filteredIcons = searchBox ? icons.filter((icon) => icon.includes(searchBox)).slice(0, 15) : [];
</script>

<InputGroup class="m-1">
  <ButtonDropdown>
    <Tooltip target={iconButton} placement="top">Change icon</Tooltip>
    <DropdownToggle color="secondary" caret bind:inner={iconButton}>
      <Icon name={preference.icon} alt="Current icon" />
    </DropdownToggle>
    <DropdownMenu>
      <p class="px-2">
        <Input type="text" placeholder="Search for icons" bind:value={searchBox} />
      </p>
      <DropdownItem divider />
      {#each filteredIcons as icon}
        <DropdownItem active={preference.icon === icon} on:click={() => (preference.icon = icon)}
          ><Icon name={icon} alt="Icon: {icon}" /> {icon}</DropdownItem
        >
      {:else}
        <p class="px-2">Start typing to filter</p>
      {/each}
    </DropdownMenu>
  </ButtonDropdown>
  <input type="text" class="form-control" bind:value={preference.tooltip} />
  <Tooltip target={sizeButton} placement="top">Change text size</Tooltip>
  <ButtonDropdown>
    <DropdownToggle color="secondary" caret bind:inner={sizeButton}>
      <Icon name="type" alt="Text size" />
    </DropdownToggle>
    <DropdownMenu>
      <DropdownItem
        active={preference.size === PreferenceSize.Large}
        on:click={() => (preference.size = PreferenceSize.Large)}>Large</DropdownItem
      >
      <DropdownItem
        active={preference.size === PreferenceSize.Normal}
        on:click={() => (preference.size = PreferenceSize.Normal)}>Medium</DropdownItem
      >
      <DropdownItem
        active={preference.size === PreferenceSize.Small}
        on:click={() => (preference.size = PreferenceSize.Small)}>Small</DropdownItem
      >
    </DropdownMenu>
  </ButtonDropdown>
  <IconButton
    color="secondary"
    icon={preference.favourite ? "star-fill" : "star"}
    click={toggleFavourite}
    active={preference.favourite}
    tooltip="Treat like favourite"
  />
  <IconButton
    color="secondary"
    icon="fonts"
    click={toggleMuted}
    active={preference.muted}
    tooltip="Show as muted text"
  />
  <IconButton color="danger" icon="trash3" tooltip="Remove preference" click={remove} />
</InputGroup>
