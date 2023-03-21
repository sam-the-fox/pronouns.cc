<script lang="ts">
  import type { Pronoun } from "$lib/api/entities";

  export let pronouns: Pronoun;

  let pronounText: string;
  $: pronounText = updatePronouns(pronouns);

  const updatePronouns = (pronouns: Pronoun) => {
    if (pronouns.display_text) {
      return pronouns.display_text;
    } else {
      const split = pronouns.pronouns.split("/");
      if (split.length < 2) return split.join("/");
      else return split.slice(0, 2).join("/");
    }
  };

  let link: string;
  let shouldLink: boolean;

  $: link = pronouns.display_text
    ? `${pronouns.pronouns},${pronouns.display_text}`
    : pronouns.pronouns;
  $: shouldLink = pronouns.pronouns.split("/").length === 5;
</script>

{#if shouldLink}
  <a href="/pronouns/{link}">{pronounText}</a>
{:else}
  {pronounText}
{/if}
