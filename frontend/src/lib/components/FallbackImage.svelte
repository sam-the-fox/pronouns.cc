<script lang="ts">
  export let urls: string[];
  export let alt: string;
  export let width = 300;

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
</script>

{#if urls.length > 0}
  <picture class="rounded-circle img-fluid">
    {#each urls as url}
      <source {width} srcSet={url} type={contentTypeFor(url)} />
    {/each}
    <img {width} src={urls[0]} {alt} class="rounded-circle img-fluid" />
  </picture>
{:else}
  <!-- TODO: actual placeholder that isn't a cat -->
  <img {width} class="rounded-circle img-fluid" src="https://placekitten.com/512/512" {alt} />
{/if}
