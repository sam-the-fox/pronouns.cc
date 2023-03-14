<script lang="ts">
  import { Card, CardHeader, CardTitle, ListGroup, ListGroupItem } from "sveltestrap";

  import type { Field } from "$lib/api/entities";

  import StatusIcon from "./StatusIcon.svelte";

  export let field: Field;
</script>

<div>
  <h3>{field.name}</h3>
  <ul class="list-unstyled fs-5">
    {#each field.entries as entry}
      <li><StatusIcon status={entry.status} /> {entry.value}</li>
    {/each}
  </ul>
</div>
