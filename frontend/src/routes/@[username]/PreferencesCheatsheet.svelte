<script lang="ts">
  import type { CustomPreferences, CustomPreference, Field, FieldEntry, Pronoun } from "$lib/api/entities";
  import defaultPreferences from "$lib/api/default_preferences";
  import StatusIcon from "$lib/components/StatusIcon.svelte";

  export let preferences: CustomPreferences;
  export let names: FieldEntry[];
  export let pronouns: Pronoun[];
  export let fields: Field[];

  let mergedPreferences: CustomPreferences;
  $: mergedPreferences = Object.assign({}, defaultPreferences, preferences);

  // Filter default preferences to the ones the user/member has used
  // This is done separately from custom preferences to make the shown list cleaner
  let usedDefaultPreferences: Array<{ id: string; preference: CustomPreference }>;
  $: usedDefaultPreferences = Object.keys(defaultPreferences)
    .filter(
      (pref) =>
        names.some((entry) => entry.status === pref) ||
        pronouns.some((entry) => entry.status === pref) ||
        fields.some((field) => field.entries.some((entry) => entry.status === pref)),
    )
    .map((key) => ({
      id: key,
      preference: defaultPreferences[key],
    }));
  // Do the same for custom preferences
  let usedCustomPreferences: Array<{ id: string; preference: CustomPreference }>;
  $: usedCustomPreferences = Object.keys(preferences)
    .filter(
      (pref) =>
        names.some((entry) => entry.status === pref) ||
        pronouns.some((entry) => entry.status === pref) ||
        fields.some((field) => field.entries.some((entry) => entry.status === pref)),
    )
    .map((pref) => ({ id: pref, preference: mergedPreferences[pref] }));
</script>

<div class="text-center">
  <ul class="list-inline text-body-secondary">
    {#each usedDefaultPreferences as pref (pref.id)}
      <li class="list-inline-item mx-2">
        <StatusIcon {preferences} status={pref.id} />
        {pref.preference.tooltip}
      </li>
    {/each}
  </ul>
  {#if usedCustomPreferences}
    <ul class="list-inline text-body-secondary">
      {#each usedCustomPreferences as pref (pref.id)}
        <li class="list-inline-item mx-2">
          <StatusIcon {preferences} status={pref.id} />
          {pref.preference.tooltip}
        </li>
      {/each}
    </ul>
  {/if}
</div>
