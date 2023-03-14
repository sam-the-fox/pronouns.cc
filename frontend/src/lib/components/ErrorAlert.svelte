<script lang="ts">
  import type { APIError } from "$lib/api/entities";
  import { Alert } from "sveltestrap";

  export let error: APIError;
</script>

<Alert color="danger" fade={false}>
  <h4 class="alert-heading">An error occurred</h4>
  <b>{error.code}:</b>
  {error.message}
  {#if error.details}
    ({error.details})
  {/if}
</Alert>
