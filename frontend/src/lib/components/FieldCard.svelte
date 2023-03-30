<script lang="ts">
  import type { Field } from "$lib/api/entities";

  import StatusLine from "./StatusLine.svelte";

  export let field: Field;
</script>

<div>
  <h3>{field.name}</h3>
  <ul class="list-unstyled fs-5">
    {#each field.entries as entry}
      <li><StatusLine status={entry.status}>{entry.value}</StatusLine></li>
    {/each}
  </ul>
</div>
