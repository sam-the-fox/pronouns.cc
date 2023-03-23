<script lang="ts">
  import type { APIError } from "$lib/api/entities";
  import { apiFetchClient } from "$lib/api/fetch";
  import ErrorAlert from "$lib/components/ErrorAlert.svelte";
  import { addToast } from "$lib/toast";
  import { Button, FormGroup, Modal, ModalBody, ModalFooter } from "sveltestrap";
  import type { PageData } from "./$types";
  import ReportCard from "./ReportCard.svelte";

  export let data: PageData;

  let warnModalOpen = false;
  const toggleWarnModal = () => (warnModalOpen = !warnModalOpen);

  let reportIndex = -1;
  let reason = "";
  let deleteUser = false;
  let error: APIError | null = null;

  const openWarnModalFor = (index: number) => {
    reportIndex = index;
    toggleWarnModal();
  };

  const warnUser = async () => {
    try {
      await apiFetchClient<any>(`/admin/reports/${data.reports[reportIndex].id}`, "PATCH", {
        warn: true,
        reason: reason,
      });
      error = null;

      addToast({ body: "Successfully warned user", header: "Warned user" });
      toggleWarnModal();
      reportIndex = -1;
    } catch (e) {
      error = e as APIError;
    }
  };
</script>

<svelte:head>
  <title>Reports - pronouns.cc</title>
</svelte:head>

<div class="container">
  <h1>Reports</h1>

  <div>
    {#each data.reports as report, index}
      <div class="my-2">
        <ReportCard {report}>
          &bull;
          <Button outline color="warning" size="sm" on:click={() => openWarnModalFor(index)}
            >Warn user</Button
          >
          <Button outline color="danger" size="sm">Deactivate user</Button>
          <Button outline color="secondary" size="sm">Ignore report</Button>
        </ReportCard>
      </div>
    {/each}
  </div>

  <Modal header="Warn user" isOpen={warnModalOpen} toggle={toggleWarnModal}>
    <ModalBody>
      {#if error}
        <ErrorAlert {error} />
      {/if}
      <ReportCard report={data.reports[reportIndex]} />
      <FormGroup floating label="Reason" class="mt-2">
        <textarea style="min-height: 100px;" class="form-control" bind:value={reason} />
      </FormGroup>
    </ModalBody>
    <ModalFooter>
      <Button color="danger" on:click={warnUser} disabled={!reason}>Warn user</Button>
      <Button color="secondary" on:click={toggleWarnModal}>Cancel</Button>
    </ModalFooter>
  </Modal>
</div>
