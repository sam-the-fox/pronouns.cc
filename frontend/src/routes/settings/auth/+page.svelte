<script lang="ts">
  import type { APIError } from "$lib/api/entities";
  import { apiFetch } from "$lib/api/fetch";
  import ErrorAlert from "$lib/components/ErrorAlert.svelte";
  import {
    Button,
    Card,
    CardBody,
    CardText,
    CardTitle,
    Input,
    Modal,
    ModalBody,
    ModalFooter,
  } from "sveltestrap";
  import type { PageData } from "./$types";

  export let data: PageData;

  let canUnlink = false;

  $: canUnlink =
    [data.user.discord, data.user.fediverse]
      .map<number>((entry) => (entry === null ? 0 : 1))
      .reduce((prev, current) => prev + current) >= 2;

  let error: APIError | null = null;
  let instance = "";
  let fediDisabled = false;

  let fediLinkModalOpen = false;
  let toggleFediLinkModal = () => (fediLinkModalOpen = !fediLinkModalOpen);

  const fediLogin = async () => {
    fediDisabled = true;
    try {
      const resp = await apiFetch<{ url: string }>(
        `/auth/urls/fediverse?instance=${encodeURIComponent(instance)}`,
        {},
      );
      window.location.assign(resp.url);
    } catch (e) {
      error = e as APIError;
    } finally {
      fediDisabled = false;
    }
  };
</script>

<div>
  <h1>Authentication providers</h1>

  <div>
    <div class="my-2">
      <Card>
        <CardBody>
          <CardTitle>Fediverse</CardTitle>
          <CardText>
            {#if data.user.fediverse}
              Your currently linked Fediverse account is <b
                >{data.user.fediverse_username}@{data.user.fediverse_instance}</b
              >
              (<code>{data.user.fediverse}</code>).
            {:else}
              You do not have a linked Fediverse account.
            {/if}
          </CardText>
          {#if data.user.fediverse}
            <Button color="danger" disabled={!canUnlink}>Unlink account</Button>
          {:else}
            <Button color="secondary" on:click={toggleFediLinkModal}>Link account</Button>
          {/if}
        </CardBody>
      </Card>
    </div>
    <div class="my-2">
      <Card>
        <CardBody>
          <CardTitle>Discord</CardTitle>
          <CardText>
            {#if data.user.discord}
              Your currently linked Discord account is <b>{data.user.discord_username}</b>
              (<code>{data.user.discord}</code>).
            {:else}
              You do not have a linked Discord account.
            {/if}
          </CardText>
          {#if data.user.discord}
            <Button color="danger" disabled={!canUnlink}>Unlink account</Button>
          {:else}
            <Button color="secondary" href={data.urls.discord}>Link account</Button>
          {/if}
        </CardBody>
      </Card>
    </div>
    <Modal header="Pick an instance" isOpen={fediLinkModalOpen} toggle={toggleFediLinkModal}>
      <ModalBody>
        <p>
          <strong>Note:</strong> Misskey (and derivatives) are not supported yet, sorry.
        </p>
        <Input placeholder="Instance (e.g. mastodon.social)" bind:value={instance} />
        {#if error}
          <div class="mt-2">
            <ErrorAlert {error} />
          </div>
        {/if}
      </ModalBody>
      <ModalFooter>
        <Button color="primary" disabled={fediDisabled || instance === ""} on:click={fediLogin}
          >Log in</Button
        >
      </ModalFooter>
    </Modal>
  </div>
</div>
