<script lang="ts">
  export let urls: string[] | null;
  export let alt: string;

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

{#if urls}
  <picture class="rounded-circle img-fluid">
    {#each urls as url}
      <source width=300 srcSet={url} type={contentTypeFor(url)} />
    {/each}
    <img width=300 src={urls[0]} {alt} class="rounded-circle img-fluid" />
  </picture>
{:else}
  <!-- TODO: actual placeholder that isn't a cat -->
  <img width=300 class="rounded-circle img-fluid" src="https://placekitten.com/512/512" {alt} />
{/if}
