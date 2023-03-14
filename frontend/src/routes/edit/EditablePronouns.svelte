<script lang="ts">
  import { WordStatus, type Pronoun } from "$lib/api/entities";
  import IconButton from "$lib/components/IconButton.svelte";
  import {
    ButtonDropdown,
    DropdownItem,
    DropdownMenu,
    DropdownToggle,
    Icon,
    Tooltip,
  } from "sveltestrap";

  export let pronoun: Pronoun;
  export let moveUp: () => void;
  export let moveDown: () => void;
  export let remove: () => void;

  let buttonElement: HTMLElement;

  const iconFor = (wordStatus: WordStatus) => {
    switch (wordStatus) {
      case WordStatus.Favourite:
        return "heart-fill";
      case WordStatus.Okay:
        return "hand-thumbs-up";
      case WordStatus.Jokingly:
        return "emoji-laughing";
      case WordStatus.FriendsOnly:
        return "people";
      case WordStatus.Avoid:
        return "hand-thumbs-down";
      default:
        return "hand-thumbs-up";
    }
  };

  const textFor = (wordStatus: WordStatus) => {
    switch (wordStatus) {
      case WordStatus.Favourite:
        return "Favourite";
      case WordStatus.Okay:
        return "Okay";
      case WordStatus.Jokingly:
        return "Jokingly";
      case WordStatus.FriendsOnly:
        return "Friends only";
      case WordStatus.Avoid:
        return "Avoid";
      default:
        return "Okay";
    }
  };
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
  <input
    placeholder="Optional display text (e.g. it/its)"
    type="text"
    class="form-control"
    bind:value={pronoun.display_text}
  />
  <ButtonDropdown>
    <Tooltip target={buttonElement} placement="top">{textFor(pronoun.status)}</Tooltip>
    <DropdownToggle color="secondary" caret bind:inner={buttonElement}>
      <Icon name={iconFor(pronoun.status)} />
    </DropdownToggle>
    <DropdownMenu>
      <DropdownItem
        on:click={() => (pronoun.status = WordStatus.Favourite)}
        active={pronoun.status === WordStatus.Favourite}>Favourite</DropdownItem
      >
      <DropdownItem
        on:click={() => (pronoun.status = WordStatus.Okay)}
        active={pronoun.status === WordStatus.Okay}>Okay</DropdownItem
      >
      <DropdownItem
        on:click={() => (pronoun.status = WordStatus.Jokingly)}
        active={pronoun.status === WordStatus.Jokingly}>Jokingly</DropdownItem
      >
      <DropdownItem
        on:click={() => (pronoun.status = WordStatus.FriendsOnly)}
        active={pronoun.status === WordStatus.FriendsOnly}>Friends only</DropdownItem
      >
      <DropdownItem
        on:click={() => (pronoun.status = WordStatus.Avoid)}
        active={pronoun.status === WordStatus.Avoid}>Avoid</DropdownItem
      >
    </DropdownMenu>
  </ButtonDropdown>
  <IconButton color="danger" icon="trash3" tooltip="Remove pronouns" click={remove} />
</div>
