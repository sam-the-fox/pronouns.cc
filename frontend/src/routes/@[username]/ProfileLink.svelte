<script lang="ts">
  import { Icon } from "sveltestrap";

  export let link: string;

  let displayLink: string;
  $: displayLink = prettifyLink(link);

  let isLink = true;
  $: isLink = link.startsWith("http://") || link.startsWith("https://");

  const prettifyLink = (raw: string) => {
    let out = raw;
    if (raw.startsWith("https://")) out = raw.substring("https://".length);
    else if (raw.startsWith("http://")) out = raw.substring("http://".length);

    if (raw.endsWith("/")) out = raw.substring(0, raw.length - 1);

    return out;
  };
</script>

<Icon name="globe" aria-hidden />
{#if isLink}
  <a href={link} rel="me nofollow noreferrer" target="_blank">{displayLink}</a>
{:else}
  <span>{link}</span>
{/if}
