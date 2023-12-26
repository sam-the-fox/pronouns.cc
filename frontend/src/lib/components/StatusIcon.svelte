<script lang="ts">
  import { Icon, Tooltip } from "@sveltestrap/sveltestrap";

  import type { CustomPreference, CustomPreferences } from "$lib/api/entities";
  import defaultPreferences from "$lib/api/default_preferences";

  export let preferences: CustomPreferences;
  export let status: string;
  export let className: string | null = null;

  let mergedPreferences: CustomPreferences;
  $: mergedPreferences = Object.assign({}, defaultPreferences, preferences);

  let currentPreference: CustomPreference;
  $: currentPreference =
    status in mergedPreferences ? mergedPreferences[status] : defaultPreferences.missing;

  let iconElement: HTMLSpanElement;
</script>

<span bind:this={iconElement} aria-hidden>
  <Icon name={currentPreference.icon} class={className} />
</span>
<span class="visually-hidden">{currentPreference.tooltip}:</span>
<Tooltip aria-hidden target={iconElement} placement="top">{currentPreference.tooltip}</Tooltip>
