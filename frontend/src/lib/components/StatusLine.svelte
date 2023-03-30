<script lang="ts">
  import { WordStatus } from "$lib/api/entities";
  import StatusIcon from "$lib/components/StatusIcon.svelte";

  export let status: WordStatus;
</script>

{#if status === WordStatus.Favourite}
  <strong class="fs-5"><StatusIcon {status} /> <slot /></strong>
{:else if status === WordStatus.Avoid}
  <span class="fs-6 text-muted"><StatusIcon {status} /> <slot /></span>
{:else}
  <StatusIcon {status} /> <slot />
{/if}
