<script lang="ts">
  import { Icon, Tooltip } from "sveltestrap";

  import { WordStatus } from "$lib/api/entities";

  export let status: WordStatus;

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

  let statusIcon: string;
  $: statusIcon = iconFor(status);

  let statusText: string;
  $: statusText = textFor(status);

  let iconElement: HTMLElement;
</script>

<span bind:this={iconElement} tabindex={0}><Icon name={statusIcon} /></span>
<Tooltip target={iconElement} placement="top">{statusText}</Tooltip>
