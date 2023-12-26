<script lang="ts">
  import type { APIError } from "$lib/api/entities";
  import { fastFetchClient } from "$lib/api/fetch";
  import ErrorAlert from "$lib/components/ErrorAlert.svelte";
  import { addToast } from "$lib/toast";
  import { Button, FormGroup, Icon, Modal, ModalBody, ModalFooter } from "@sveltestrap/sveltestrap";

  export let subject: string;
  export let reportUrl: string;

  const MAX_REASON_LENGTH = 2000;
  let reason = "";

  let error: APIError | null = null;

  let isOpen = false;
  const toggle = () => (isOpen = !isOpen);

  const doReport = async () => {
    try {
      await fastFetchClient(reportUrl, "POST", { reason: reason });
      error = null;
      addToast({ header: "Sent report", body: "Successfully sent report!" });
      toggle();
    } catch (e) {
      error = e as APIError;
    }
  };
</script>

<Button color="danger" outline on:click={toggle}
  ><Icon name="exclamation-triangle-fill" /> Report {subject}</Button
>

<Modal header="Report {subject}" {isOpen} {toggle}>
  <ModalBody>
    {#if error}
      <ErrorAlert {error} />
    {/if}
    <FormGroup floating label="Reason" class="mt-2">
      <textarea style="min-height: 100px;" class="form-control" bind:value={reason} />
    </FormGroup>
    <p class="text-muted mt-3">
      <Icon name="info-circle-fill" aria-hidden /> Please only report {subject}s for violating the
      <a class="text-reset" href="/page/terms" target="_blank" rel="noopener noreferrer"
        >terms of service</a
      >.
      <br />
      Note that we <strong>cannot</strong> take action on any reports that cannot be proven with just
      the reported profile and any external pages linked on it.
    </p>
  </ModalBody>
  <ModalFooter>
    <Button
      color="danger"
      disabled={!reason || reason.length > MAX_REASON_LENGTH}
      on:click={doReport}>Report</Button
    >
    <Button color="secondary" on:click={toggle}>Cancel</Button>
  </ModalFooter>
</Modal>
