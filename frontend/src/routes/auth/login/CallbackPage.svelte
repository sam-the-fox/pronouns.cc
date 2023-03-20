<script lang="ts">
  import { goto } from "$app/navigation";
  import type { APIError, MeUser } from "$lib/api/entities";
  import { apiFetch } from "$lib/api/fetch";
  import ErrorAlert from "$lib/components/ErrorAlert.svelte";
  import { userStore } from "$lib/store";
  import { addToast } from "$lib/toast";
  import { onMount } from "svelte";
  import {
    Alert,
    Button,
    FormGroup,
    Icon,
    Input,
    Modal,
    ModalBody,
    ModalFooter,
  } from "sveltestrap";

  export let authType: string;
  export let remoteName: string | undefined;
  export let error: APIError | undefined;
  export let requireInvite: boolean | undefined;
  export let isDeleted: boolean | undefined;
  export let ticket: string | undefined;
  export let token: string | undefined;
  export let user: MeUser | undefined;
  export let deletedAt: string | undefined;

  onMount(() => {
    if (!isDeleted && token && user) {
      localStorage.setItem("pronouns-token", token);
      localStorage.setItem("pronouns-user", JSON.stringify(user));
      userStore.set(user);
      goto("/");
    }
  });

  let deleteCancelled: boolean;
  let deleteError: APIError | null;

  let username: string;
  let inviteCode: string;
  let forceDeleteName = "";
  let forceDeleteModalOpen = false;
  let toggleForceDeleteModal = () => (forceDeleteModalOpen = !forceDeleteModalOpen);

  export let linkAccount: () => Promise<void>;
  export let signupForm: (username: string, inviteCode: string) => Promise<void>;

  const forceDeleteAccount = async () => {
    try {
      await apiFetch<any>("/auth/force-delete", {
        method: "GET",
        headers: {
          "X-Delete-Token": token!,
        },
      });

      deleteError = null;
      addToast({ header: "Deleted account", body: "Successfully deleted your account" });
      goto("/");
    } catch (e) {
      deleteError = e as APIError;
    }
  };

  const cancelDelete = async () => {
    try {
      await apiFetch<any>("/auth/cancel-delete", {
        method: "GET",
        headers: {
          "X-Delete-Token": token!,
        },
      });

      deleteCancelled = true;
      deleteError = null;
    } catch (e) {
      deleteCancelled = false;
      deleteError = e as APIError;
    }
  };
</script>

<svelte:head>
  <title>Log in with the {authType} - pronouns.cc</title>
</svelte:head>

<h1>Log in with the {authType}</h1>

{#if error}
  <ErrorAlert {error} />
{/if}
{#if ticket && $userStore}
  <div>
    <FormGroup floating label="{authType} username">
      <Input readonly value={remoteName} />
    </FormGroup>
  </div>
  <div class="my-2">
    <FormGroup floating label="pronouns.cc username">
      <Input readonly value={$userStore.name} />
    </FormGroup>
  </div>
  <div>
    <Button on:click={linkAccount}>Link account</Button>
    <Button color="secondary" href="/settings/auth">Cancel</Button>
  </div>
{:else if ticket}
  <form on:submit|preventDefault={() => signupForm(username, inviteCode)}>
    <div>
      <FormGroup floating label="{authType} username">
        <Input readonly value={remoteName} />
      </FormGroup>
    </div>
    <div>
      <FormGroup floating label="Username">
        <Input id="username" name="username" bind:value={username} />
      </FormGroup>
    </div>
    {#if requireInvite}
      <div>
        <FormGroup floating label="Invite code">
          <Input id="invite" name="invite" aria-describedby="invite-help" bind:value={inviteCode} />
        </FormGroup>
        <div id="invite-help" class="form-text">
          <Icon name="info-circle-fill" /> You currently need an invite code to sign up. You can get
          one from an existing user.
        </div>
      </div>
    {/if}
    <div class="form-text mb-1">
      By signing up, you agree to the <a href="/page/tos">terms of service</a> and the
      <a href="/page/privacy">privacy policy</a>.
    </div>
    <Button type="submit" color="primary">Sign up</Button>
  </form>
{:else if isDeleted && token}
  <p>Your account is pending deletion since {deletedAt}.</p>
  <p>If you wish to cancel deletion, press the button below.</p>
  <p>
    <Button color="primary" on:click={cancelDelete} disabled={deleteCancelled}
      >Cancel account deletion</Button
    >
  </p>
  <p>
    Alternatively, if you want your data wiped immediately, press the force delete link below. <b
      >This is irreversible.</b
    >
  </p>
  <p>
    <Button color="link" on:click={toggleForceDeleteModal}>Force delete account</Button>
  </p>
  <Modal
    header="Force delete account"
    isOpen={forceDeleteModalOpen}
    toggle={toggleForceDeleteModal}
  >
    <ModalBody>
      <p>
        If you want to delete your account, type your username below:
        <br />
        <b>
          This is irreversible! Your account <i>cannot</i> be recovered after you press "Force delete
          account".
        </b>
      </p>
      <p>
        <input type="text" class="form-control" bind:value={forceDeleteName} />
      </p>
      {#if deleteError}
        <ErrorAlert error={deleteError} />
      {/if}
    </ModalBody>
    <ModalFooter>
      <Button color="danger" on:click={forceDeleteAccount} disabled={forceDeleteName !== user?.name}
        >Force delete account</Button
      >
      <Button color="secondary" on:click={toggleForceDeleteModal}>Cancel delete</Button>
    </ModalFooter>
  </Modal>
  {#if deleteCancelled}
    <Alert color="secondary" fade={false}>
      Account deletion cancelled! You can now <a href="/auth/login">log in</a> again.
    </Alert>
  {/if}
  {#if deleteError}
    <ErrorAlert error={deleteError} />
  {/if}
{:else}
  Loading...
{/if}
