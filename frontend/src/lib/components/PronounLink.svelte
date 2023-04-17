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
      if (split.length === 5) return split.splice(0, 2).join("/");
      return pronouns.pronouns;
    }
  };

  let link: string;
  let shouldLink: boolean;

  $: link = linkPronouns(pronouns);
  $: shouldLink = pronouns.pronouns.split("/").length === 5;

  const linkPronouns = (pronouns: Pronoun) => {
    const linkBase = pronouns.pronouns
      .split("/")
      .map((snippet) => encodeURIComponent(snippet))
      .join("/");

    if (pronouns.display_text) {
      return `${linkBase},${encodeURIComponent(pronouns.display_text)}`;
    }
    return linkBase;
  };
</script>

{#if shouldLink}
  <a class="text-reset" href="/pronouns/{link}">{pronounText}</a>
{:else}
  {pronounText}
{/if}
