<script lang="ts">
  import { defaultAvatars } from "$lib/api/entities";

  export let urls: string[];
  export let alt: string;
  export let width = 300;
  export let lazyLoad = false;

  const contentTypeFor = (url: string) => {
    if (url.endsWith(".webp")) {
      return "image/webp";
    } else if (url.endsWith(".jpg") || url.endsWith(".jpeg")) {
      return "image/jpeg";
    } else if (url.endsWith(".png")) {
      return "image/png";
    } else {
      return "application/octet-stream";
    }
  };

  $: urls = urls.length > 0 ? urls : defaultAvatars;
</script>

{#if urls.length > 0}
  <picture class="rounded-circle img-fluid">
    {#each urls as url}
      <source {width} srcSet={url} type={contentTypeFor(url)} />
    {/each}
    <img
      {width}
      height={width}
      src={urls[0] || defaultAvatars[0]}
      {alt}
      class="rounded-circle img-fluid"
      loading={lazyLoad ? "lazy" : "eager"}
    />
  </picture>
{:else}
  <img {width} height={width} class="rounded-circle img-fluid" src={defaultAvatars[0]} {alt} />
{/if}
