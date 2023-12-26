<script lang="ts">
  import { PreferenceSize, type CustomPreference } from "$lib/api/entities";
  import IconButton from "$lib/components/IconButton.svelte";
  import {
    ButtonDropdown,
    DropdownItem,
    DropdownMenu,
    DropdownToggle,
    Icon,
    InputGroup,
    Tooltip,
  } from "@sveltestrap/sveltestrap";
  import IconPicker from "./IconPicker.svelte";

  export let preference: CustomPreference;
  export let remove: () => void;

  let sizeButton: HTMLElement;

  const toggleMuted = () => (preference.muted = !preference.muted);
  const toggleFavourite = () => (preference.favourite = !preference.favourite);
</script>

<InputGroup class="m-1">
  <IconPicker bind:icon={preference.icon} />
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
