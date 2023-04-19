<script lang="ts">
  import { PreferenceSize } from "$lib/api/entities";
  import StatusIcon from "$lib/components/StatusIcon.svelte";

  import type { CustomPreference, CustomPreferences } from "$lib/api/entities";
  import defaultPreferences from "$lib/api/default_preferences";

  export let preferences: CustomPreferences;
  export let status: string;

  let mergedPreferences: CustomPreferences;
  $: mergedPreferences = Object.assign(defaultPreferences, preferences);

  let currentPreference: CustomPreference;
  $: currentPreference =
    status in mergedPreferences ? mergedPreferences[status] : defaultPreferences.okay;

  let classes: string;
  $: classes = setClasses(currentPreference);

  const setClasses = (pref: CustomPreference) => {
    let classes = "";
    if (pref.muted) {
      classes += "text-muted ";
    }
    switch (pref.size) {
      case PreferenceSize.Large:
        classes += "fs-5";
        break;
      case PreferenceSize.Normal:
        break;
      case PreferenceSize.Small:
        classes += "fs-6";
    }

    return classes.trim();
  };
</script>

{#if currentPreference.size === PreferenceSize.Large}
  <strong class={classes}><StatusIcon {preferences} {status} /> <slot /></strong>
{:else}
  <span class={classes}><StatusIcon {preferences} {status} /> <slot /></span>
{/if}
