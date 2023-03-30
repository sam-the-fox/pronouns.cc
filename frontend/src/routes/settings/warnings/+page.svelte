<script lang="ts">
  import type { APIError } from "$lib/api/entities";
  import { apiFetchClient, fastFetchClient } from "$lib/api/fetch";
  import ErrorAlert from "$lib/components/ErrorAlert.svelte";
  import { addToast } from "$lib/toast";
  import { DateTime } from "luxon";
  import { Button, Card, CardBody, CardFooter, CardHeader } from "sveltestrap";
  import type { PageData } from "./$types";

  export let data: PageData;
  let error: APIError | null = null;

  const acknowledgeWarning = async (idx: number) => {
    try {
      await fastFetchClient(`/auth/warnings/${data.warnings[idx].id}/ack`, "POST");
      addToast({
        header: "Acknowledged",
        body: `Marked warning #${data.warnings[idx].id} as read.`,
      });
      data.warnings[idx].read = true;
      data.warnings = data.warnings;
    } catch (e) {
      error = e as APIError;
    }
  };
</script>

<h1>Warnings ({data.warnings.length})</h1>

{#if error}
  <ErrorAlert {error} />
{/if}

<div>
  {#each data.warnings as warning, index}
    <Card class="my-2">
      <CardHeader>
        <strong>#{warning.id}</strong> ({DateTime.fromISO(warning.created_at)
          .toLocal()
          .toLocaleString(DateTime.DATETIME_MED)})
      </CardHeader>
      <CardBody>
        <blockquote class="blockquote">{warning.reason}</blockquote>
      </CardBody>
      {#if !warning.read}
        <CardFooter>
          <Button color="secondary" outline on:click={() => acknowledgeWarning(index)}
            >Mark as read</Button
          >
        </CardFooter>
      {/if}
    </Card>
  {:else}
    You have no warnings!
  {/each}
</div>
