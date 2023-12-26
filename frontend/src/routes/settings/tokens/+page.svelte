<script lang="ts">
  import type { APIError } from "$lib/api/entities";
  import { apiFetchClient } from "$lib/api/fetch";
  import { DateTime } from "luxon";
  import {
    Button,
    ButtonGroup,
    Card,
    CardBody,
    CardHeader,
    Modal,
    ModalBody,
    ModalFooter,
    ModalHeader,
  } from "@sveltestrap/sveltestrap";
  import type { PageData } from "./$types";

  interface Token {
    token: string;
    id: string;
    api_only: boolean;
    read_only: boolean;
    created: string;
    expires: string;
  }

  export let data: PageData;

  let newToken: string | null = null;
  let error: APIError | null = null;

  let newTokenModalOpen = false;
  const toggleNewTokenModalOpen = () => (newTokenModalOpen = !newTokenModalOpen);

  const createToken = async (readOnly: boolean) => {
    try {
      const token = await apiFetchClient<Token>(`/auth/tokens?read_only=${readOnly}`, "POST");
      newToken = token.token;
      error = null;
      data.tokens = [...data.tokens, token];
      toggleNewTokenModalOpen();
    } catch (e) {
      error = e as APIError;
    }
  };
</script>

<h1>
  API tokens ({data.tokens.length})
  <ButtonGroup>
    <Button color="success" on:click={() => createToken(false)}>New API token</Button>
    <Button color="success" on:click={() => createToken(true)}>New read-only API token</Button>
    <Button color="secondary" href="https://docs.pronouns.cc/api/" target="_blank">API docs</Button>
  </ButtonGroup>
</h1>

{#each data.tokens as token}
  <Card class="my-2">
    <CardHeader>{token.id}</CardHeader>
    <CardBody>
      <ul class="list-unstyled">
        <li>
          <strong>Created at:</strong>
          {DateTime.fromISO(token.created).toLocal().toLocaleString(DateTime.DATETIME_MED)}
        </li>
        <li>
          <strong>Expires at:</strong>
          {DateTime.fromISO(token.expires).toLocal().toLocaleString(DateTime.DATETIME_MED)}
        </li>
        <li>
          <strong>Read-only:</strong>
          {token.read_only ? "yes" : "no"}
        </li>
      </ul>
    </CardBody>
  </Card>
{:else}
  You don't have any unexpired API tokens right now.
{/each}

<Modal isOpen={newTokenModalOpen}>
  <ModalHeader toggle={toggleNewTokenModalOpen}>New token created</ModalHeader>
  <ModalBody>
    <p>Created a new API token! Please save it somewhere secure, as it will only be shown once.</p>
    <p><code>{newToken}</code></p>
  </ModalBody>
  <ModalFooter>
    <Button color="secondary" on:click={toggleNewTokenModalOpen}>Close</Button>
  </ModalFooter>
</Modal>
