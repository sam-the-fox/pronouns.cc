<script lang="ts">
  import defaultPreferences from "$lib/api/default_preferences";
  import type { CustomPreference, CustomPreferences } from "$lib/api/entities";
  import IconButton from "$lib/components/IconButton.svelte";
  import {
    ButtonDropdown,
    DropdownItem,
    DropdownMenu,
    DropdownToggle,
    Icon,
    Tooltip,
  } from "@sveltestrap/sveltestrap";

  export let value: string;
  export let status: string;
  export let preferences: CustomPreferences;
  export let moveUp: () => void;
  export let moveDown: () => void;
  export let remove: () => void;

  let buttonElement: HTMLElement;

  let mergedPreferences: CustomPreferences;
  $: mergedPreferences = Object.assign({}, defaultPreferences, preferences);

  let currentPreference: CustomPreference;
  $: currentPreference =
    status in mergedPreferences ? mergedPreferences[status] : defaultPreferences.missing;

  let preferenceIds: string[];
  $: preferenceIds = Object.keys(mergedPreferences).filter((s) => s !== "missing");
</script>

<div class="input-group m-1">
  <IconButton icon="chevron-up" color="secondary" tooltip="Move name up" click={moveUp} />
  <IconButton icon="chevron-down" color="secondary" tooltip="Move name down" click={moveDown} />
  <input type="text" class="form-control" bind:value />
  <ButtonDropdown>
    <Tooltip target={buttonElement} placement="top">{currentPreference.tooltip}</Tooltip>
    <DropdownToggle color="secondary" caret bind:inner={buttonElement}>
      <Icon name={currentPreference.icon} />
    </DropdownToggle>
    <DropdownMenu>
      {#each preferenceIds as id}
        <DropdownItem on:click={() => (status = id)} active={status === id}>
          <Icon name={mergedPreferences[id].icon} aria-hidden />
          {mergedPreferences[id].tooltip}
        </DropdownItem>
      {/each}
    </DropdownMenu>
  </ButtonDropdown>
  <IconButton color="danger" icon="trash3" tooltip="Remove name" click={remove} />
</div>
