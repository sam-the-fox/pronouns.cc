<script lang="ts">
  import defaultPreferences from "$lib/api/default_preferences";
  import type { Pronoun, CustomPreference, CustomPreferences } from "$lib/api/entities";
  import IconButton from "$lib/components/IconButton.svelte";
  import {
    Button,
    ButtonDropdown,
    Collapse,
    DropdownItem,
    DropdownMenu,
    DropdownToggle,
    Icon,
    InputGroupText,
    Popover,
    Tooltip,
  } from "sveltestrap";

  export let pronoun: Pronoun;
  export let preferences: CustomPreferences;
  export let moveUp: () => void;
  export let moveDown: () => void;
  export let remove: () => void;

  let buttonElement: HTMLElement;
  let displayOpen = false;
  const toggleDisplay = () => (displayOpen = !displayOpen);

  let mergedPreferences: CustomPreferences;
  $: mergedPreferences = Object.assign(defaultPreferences, preferences);

  let currentPreference: CustomPreference;
  $: currentPreference =
    pronoun.status in mergedPreferences
      ? mergedPreferences[pronoun.status]
      : defaultPreferences.missing;

  let preferenceIds: string[];
  $: preferenceIds = Object.keys(mergedPreferences).filter((s) => s !== "missing");
</script>

<div class="input-group m-1">
  <IconButton icon="chevron-up" color="secondary" tooltip="Move pronouns up" click={moveUp} />
  <IconButton icon="chevron-down" color="secondary" tooltip="Move pronouns down" click={moveDown} />
  <input
    placeholder="Full set (e.g. it/it/its/its/itself)"
    type="text"
    class="form-control"
    bind:value={pronoun.pronouns}
    required
  />
  <IconButton
    icon={pronoun.display_text ? "tag-fill" : "tag"}
    color="secondary"
    tooltip={pronoun.display_text ? "Edit or remove display text" : "Add display text"}
    click={toggleDisplay}
  />
  <ButtonDropdown>
    <Tooltip target={buttonElement} placement="top">{currentPreference.tooltip}</Tooltip>
    <DropdownToggle color="secondary" caret bind:inner={buttonElement}>
      <Icon name={currentPreference.icon} />
    </DropdownToggle>
    <DropdownMenu>
      {#each preferenceIds as id}
        <DropdownItem on:click={() => (pronoun.status = id)} active={pronoun.status === id}>
          <Icon name={mergedPreferences[id].icon} aria-hidden />
          {mergedPreferences[id].tooltip}
        </DropdownItem>
      {/each}
    </DropdownMenu>
  </ButtonDropdown>
  <IconButton color="danger" icon="trash3" tooltip="Remove pronouns" click={remove} />
</div>
<Collapse isOpen={displayOpen}>
  <div class="input-group">
    <InputGroupText>Display text</InputGroupText>
    <input
      placeholder="Optional display text (e.g. it/its)"
      type="text"
      class="form-control"
      bind:value={pronoun.display_text}
    />
    <IconButton id="display-help" icon="question" tooltip="Help" color="secondary" />
    <Popover target="display-help" placement="bottom">
      This is the short text shown on your profile page. If you leave it empty, it will default to
      the first two forms of the full set.
    </Popover>
  </div>
</Collapse>
