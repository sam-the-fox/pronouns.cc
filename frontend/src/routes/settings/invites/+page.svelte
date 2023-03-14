<script lang="ts">
  import type { APIError, Invite } from "$lib/api/entities";
  import { apiFetchClient } from "$lib/api/fetch";
  import ErrorAlert from "$lib/components/ErrorAlert.svelte";
  import { Alert, Button, Modal, Table } from "sveltestrap";
  import type { PageData } from "./$types";

  export let data: PageData;

  let error: APIError | null = null;
  let latestInvite: Invite | null = null;
  let open = false;

  const toggle = () => (open = !open);
  const createInvite = async () => {
    try {
      const invite = await apiFetchClient<Invite>("/auth/invites", "POST");

      error = null;
      data.invites = [...data.invites, invite];
      latestInvite = invite;
      open = true;
    } catch (e) {
      latestInvite = null;
      error = e as APIError;
    }
  };
</script>

<h1>Invites ({data.invites.length})</h1>

<div>
  {#if !data.invitesEnabled}
    <p>Invites aren't required to sign up to pronouns.cc right now!</p>
  {:else}
    <Table striped hover>
      <thead>
        <th>Code</th>
        <th>Created at</th>
        <th>Used?</th>
      </thead>
      <tbody>
        {#each data.invites as invite}
          <tr>
            <td><code>{invite.code}</code></td>
            <td>{invite.created}</td>
            <td>{invite.used ? "yes" : "no"}</td>
          </tr>
        {/each}
      </tbody>
    </Table>
    <div class="row">
      <div class="col-md-4">
        <Button color="primary" on:click={createInvite}>Create invite</Button>
      </div>
      <div class="col-md">
        {#if error}
          <ErrorAlert {error} />
        {/if}
      </div>
    </div>
    <Modal body header="Invite created" isOpen={open} {toggle}>
      Successfully created a new invite! Give the person you're sending it to this code:
      <code>{latestInvite?.code}</code>
    </Modal>
  {/if}
</div>
