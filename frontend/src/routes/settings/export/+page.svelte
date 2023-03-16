<script lang="ts">
  import type { PageData } from "./$types";
  import { DateTime, Duration } from "luxon";
  import { Alert, Button } from "sveltestrap";
  import { PUBLIC_BASE_URL } from "$env/static/public";
  import { apiFetchClient } from "$lib/api/fetch";
  import type { APIError } from "$lib/api/entities";
  import { addToast } from "$lib/toast";
  import ErrorAlert from "$lib/components/ErrorAlert.svelte";

  export let data: PageData;

  const createdAt = data.exportData?.created_at
    ? DateTime.fromISO(data.exportData.created_at).toLocal()
    : DateTime.fromSeconds(0);
  const now = DateTime.now().toLocal();

  const availableFor = Duration.fromObject({ days: 7 });
  const minTimeBetween = Duration.fromObject({ days: 1 });
  const durationSinceCreated = now.diff(createdAt, "days");

  let error: APIError | null = null;

  const requestExport = async () => {
    try {
      await apiFetchClient<any>("/users/@me/export/start");

      addToast({
        header: "Export in progress",
        body: "Your data is now being exported! Check back here in a few minutes.",
      });
    } catch (e) {
      error = e as APIError;
    }
  };
</script>

<div>
  <h1>
    Data export
    <Button
      color="success"
      disabled={durationSinceCreated.days < minTimeBetween.days}
      on:click={requestExport}>Export your data</Button
    >
  </h1>

  {#if error}
    <ErrorAlert {error} />
  {/if}

  {#if data.exportData}
    <div>
      {#if durationSinceCreated.days < minTimeBetween.days}
        <Alert color="secondary" fade={false}>
          You can only export your data once a day. You can next export your data at <b
            >{createdAt.plus(minTimeBetween).toLocaleString(DateTime.DATETIME_MED)}</b
          >.
        </Alert>
      {/if}
      <p>
        You last exported your data at {createdAt.toLocaleString(DateTime.DATETIME_MED)}.
        <br />
        This file will be available until {createdAt
          .plus(availableFor)
          .toLocal()
          .toLocaleString(DateTime.DATETIME_MED)}
        <br />
        Download your export file below:
      </p>
      <p>
        <Button color="primary" href="{PUBLIC_BASE_URL}/media{data.exportData.path}" target="_blank"
          >Download data export</Button
        >
      </p>
    </div>
  {:else}
    <Alert color="secondary" fade={false}>
      You haven't exported your data in the past week. To create an export file, click the button
      above.
    </Alert>
  {/if}
</div>
