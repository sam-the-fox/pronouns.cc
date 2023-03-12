<script lang="ts">
  import { Card, CardHeader, CardTitle, ListGroup, ListGroupItem } from "sveltestrap";

  import type { Field } from "$lib/api/entities";

  import StatusIcon from "./StatusIcon.svelte";

  export let field: Field;
</script>

<div>
  <h4>{field.name}</h4>
  <ul class="list-unstyled">
    {#each field.entries as entry}
      <li><StatusIcon status={entry.status} /> {entry.value}</li>
    {/each}
  </ul>
</div>
