<script lang="ts">
  import type { Pronoun } from "$lib/api/entities";

  export let pronouns: Pronoun;

  let pronounText: string;
  if (pronouns.display_text) {
    pronounText = pronouns.display_text;
  } else {
    const split = pronouns.pronouns.split("/");
    if (split.length < 2) pronounText = split.join("/");
    else pronounText = split.slice(0, 2).join("/");
  }

  const shouldLink = pronouns.pronouns.split("/").length === 5;
</script>

{#if shouldLink}
  <a href="/pronouns/{pronouns.pronouns}">{pronounText}</a>
{:else}
  {pronounText}
{/if}
