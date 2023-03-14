<script lang="ts">
  import { WordStatus } from "$lib/api/entities";
  import IconButton from "$lib/components/IconButton.svelte";
  import {
    ButtonDropdown,
    DropdownItem,
    DropdownMenu,
    DropdownToggle,
    Icon,
    Tooltip,
  } from "sveltestrap";

  export let value: string;
  export let status: WordStatus;
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
  <IconButton icon="chevron-up" color="secondary" tooltip="Move name up" click={moveUp} />
  <IconButton icon="chevron-down" color="secondary" tooltip="Move name down" click={moveDown} />
  <input type="text" class="form-control" bind:value />
  <ButtonDropdown>
    <Tooltip target={buttonElement} placement="top">{textFor(status)}</Tooltip>
    <DropdownToggle color="secondary" caret bind:inner={buttonElement}>
      <Icon name={iconFor(status)} />
    </DropdownToggle>
    <DropdownMenu>
      <DropdownItem
        on:click={() => (status = WordStatus.Favourite)}
        active={status === WordStatus.Favourite}>Favourite</DropdownItem
      >
      <DropdownItem on:click={() => (status = WordStatus.Okay)} active={status === WordStatus.Okay}
        >Okay</DropdownItem
      >
      <DropdownItem
        on:click={() => (status = WordStatus.Jokingly)}
        active={status === WordStatus.Jokingly}>Jokingly</DropdownItem
      >
      <DropdownItem
        on:click={() => (status = WordStatus.FriendsOnly)}
        active={status === WordStatus.FriendsOnly}>Friends only</DropdownItem
      >
      <DropdownItem
        on:click={() => (status = WordStatus.Avoid)}
        active={status === WordStatus.Avoid}>Avoid</DropdownItem
      >
    </DropdownMenu>
  </ButtonDropdown>
  <IconButton color="danger" icon="trash3" tooltip="Remove name" click={remove} />
</div>
