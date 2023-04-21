<script lang="ts">
  import defaultPreferences from "$lib/api/default_preferences";
  import {
    memberAvatars,
    type PartialMember,
    type User,
    type CustomPreferences,
  } from "$lib/api/entities";
  import FallbackImage from "./FallbackImage.svelte";

  export let user: User;
  export let member: PartialMember;

  let pronouns: string | undefined;

  let mergedPreferences: CustomPreferences;
  $: mergedPreferences = Object.assign(defaultPreferences, user.custom_preferences);

  const getPronouns = (member: PartialMember) => {
    const filteredPronouns = member.pronouns.filter(
      (entry) => (mergedPreferences[entry.status] || { favourite: false }).favourite,
    );
    if (filteredPronouns.length === 0) {
      return undefined;
    }
    return filteredPronouns
      .map((pronouns) => {
        if (pronouns.display_text) {
          return pronouns.display_text;
        } else {
          const split = pronouns.pronouns.split("/");
          if (split.length === 5) return split.splice(0, 2).join("/");
          return pronouns.pronouns;
        }
      })
      .join(", ");
  };

  $: pronouns = getPronouns(member);
</script>

<div>
  <a href="/@{user.name}/{member.name}">
    <FallbackImage urls={memberAvatars(member)} width={200} alt="Avatar for {member.name}" />
  </a>
  <p class="m-2">
    <a class="text-reset fs-5 text-break" href="/@{user.name}/{member.name}">
      {member.display_name ?? member.name}
    </a>
    {#if pronouns}
      <br />
      {pronouns}
    {/if}
  </p>
</div>
