<script lang="ts">
  import type { APIError } from "$lib/api/entities";
  import { fastFetchClient } from "$lib/api/fetch";
  import ErrorAlert from "$lib/components/ErrorAlert.svelte";
  import { addToast } from "$lib/toast";
  import { Button, ButtonGroup, FormGroup, Modal, ModalBody, ModalFooter } from "sveltestrap";
  import type { PageData } from "./$types";
  import ReportCard from "./ReportCard.svelte";

  export let data: PageData;

  let warnModalOpen = false;
  const toggleWarnModal = () => (warnModalOpen = !warnModalOpen);

  let banModalOpen = false;
  const toggleBanModal = () => (banModalOpen = !banModalOpen);

  let ignoreModalOpen = false;
  const toggleIgnoreModal = () => (ignoreModalOpen = !ignoreModalOpen);

  let reportIndex = -1;
  let reason = "";
  let deleteUser = false;
  let error: APIError | null = null;

  $: console.log(deleteUser);

  const openWarnModalFor = (index: number) => {
    reportIndex = index;
    toggleWarnModal();
  };

  const openBanModalFor = (index: number) => {
    reportIndex = index;
    toggleBanModal();
  };

  const openIgnoreModalFor = (index: number) => {
    reportIndex = index;
    toggleIgnoreModal();
  };

  const warnUser = async () => {
    try {
      await fastFetchClient(`/admin/reports/${data.reports[reportIndex].id}`, "PATCH", {
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

  const deactivateUser = async () => {
    try {
      await fastFetchClient(`/admin/reports/${data.reports[reportIndex].id}`, "PATCH", {
        warn: true,
        ban: true,
        delete: deleteUser,
        reason: reason,
      });
      error = null;

      addToast({ body: "Successfully deactivated user", header: "Deactivated user" });
      toggleBanModal();
      reportIndex = -1;
    } catch (e) {
      error = e as APIError;
    }
  };

  const ignoreReport = async () => {
    try {
      await fastFetchClient(`/admin/reports/${data.reports[reportIndex].id}`, "PATCH", {
        reason: reason,
      });
      error = null;

      addToast({ body: "Successfully acknowledged report", header: "Ignored report" });
      toggleIgnoreModal();
      reportIndex = -1;
    } catch (e) {
      error = e as APIError;
    }
  };

  const urlParamsWith = (name: string, value: string) => {
    const params = new URLSearchParams(window.location.search);
    params.set(name, value);
    return params.toString();
  };

  const urlParamsWithout = (name: string) => {
    const params = new URLSearchParams(window.location.search);
    params.delete(name);
    return params.toString();
  };
</script>

<svelte:head>
  <title>Reports - pronouns.cc</title>
</svelte:head>

<div class="container">
  <h1>Reports</h1>

  <div>
    <ButtonGroup>
      {#if data.userId || data.reporterId}
        <Button color="secondary" disabled>Showing all reports</Button>
      {:else if data.isClosed}
        <Button color="primary" href="/reports?{urlParamsWith('closed', 'false')}"
          >Showing closed reports</Button
        >
      {:else}
        <Button color="secondary" href="/reports?{urlParamsWith('closed', 'true')}"
          >Showing open reports</Button
        >
      {/if}
      {#if data.userId}
        <Button color="secondary" href="/reports?{urlParamsWithout('user_id')}"
          >Filtering on specific user</Button
        >
      {:else}
        <Button color="secondary" disabled>Not filtering on user</Button>
      {/if}
      {#if data.reporterId}
        <Button color="secondary" href="/reports?{urlParamsWithout('reporter_id')}"
          >Filtering on specific reporter</Button
        >
      {:else}
        <Button color="secondary" disabled>Not filtering on reporter</Button>
      {/if}
    </ButtonGroup>
  </div>

  <div>
    {#each data.reports as report, index}
      <div class="my-2">
        <ReportCard {report}>
          &bull;
          <Button outline color="warning" size="sm" on:click={() => openWarnModalFor(index)}
            >Warn user</Button
          >
          <Button outline color="danger" size="sm" on:click={() => openBanModalFor(index)}
            >Deactivate user</Button
          >
          <Button outline color="secondary" size="sm" on:click={() => openIgnoreModalFor(index)}
            >Ignore report</Button
          >
          &bull;
          <a href="/reports?user_id={report.user_id}">Show all reports of this user</a>
          &bull;
          <a href="/reports?reporter_id={report.reporter_id}">Show all reports by this reporter</a>
        </ReportCard>
      </div>
    {:else}
      There are no reports matching your search!
    {/each}
  </div>

  {#if data.reports.length >= 100}
    <div>
      <a
        href="/reports?{urlParamsWith(
          'before',
          data.reports[data.reports.length - 1].id.toString(),
        )}">Load more reports</a
      >
    </div>
  {/if}

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

  <Modal header="Deactivate user" isOpen={banModalOpen} toggle={toggleBanModal}>
    <ModalBody>
      {#if error}
        <ErrorAlert {error} />
      {/if}
      <ReportCard report={data.reports[reportIndex]} />
      <FormGroup floating label="Reason" class="my-2">
        <textarea style="min-height: 100px;" class="form-control" bind:value={reason} />
      </FormGroup>
      <div class="form-check">
        <input class="form-check-input" type="checkbox" bind:checked={deleteUser} id="deleteUser" />
        <label class="form-check-label" for="deleteUser">Delete user?</label>
      </div>
    </ModalBody>
    <ModalFooter>
      <Button color="danger" on:click={deactivateUser} disabled={!reason}>Deactivate user</Button>
      <Button color="secondary" on:click={toggleBanModal}>Cancel</Button>
    </ModalFooter>
  </Modal>

  <Modal header="Ignore report" isOpen={ignoreModalOpen} toggle={toggleIgnoreModal}>
    <ModalBody>
      {#if error}
        <ErrorAlert {error} />
      {/if}
      <ReportCard report={data.reports[reportIndex]} />
      <FormGroup floating label="Reason" class="my-2">
        <textarea style="min-height: 100px;" class="form-control" bind:value={reason} />
      </FormGroup>
    </ModalBody>
    <ModalFooter>
      <Button color="warning" on:click={ignoreReport} disabled={!reason}>Ignore report</Button>
      <Button color="secondary" on:click={toggleIgnoreModal}>Cancel</Button>
    </ModalFooter>
  </Modal>
</div>
