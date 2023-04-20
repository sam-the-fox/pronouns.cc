<script lang="ts">
  import { goto } from "$app/navigation";
  import { type MeUser, userAvatars, type APIError, MAX_MEMBERS } from "$lib/api/entities";
  import { apiFetchClient, fastFetchClient } from "$lib/api/fetch";
  import { usernameRegex } from "$lib/api/regex";
  import ErrorAlert from "$lib/components/ErrorAlert.svelte";
  import FallbackImage from "$lib/components/FallbackImage.svelte";
  import { userStore } from "$lib/store";
  import { addToast } from "$lib/toast";
  import {
    Alert,
    Button,
    Icon,
    Modal,
    ModalBody,
    ModalFooter,
    ModalHeader,
    Table,
  } from "sveltestrap";
  import type { PageData } from "./$types";
  import { onMount } from "svelte";

  export let data: PageData;

  let username = data.user.name;
  let usernameValid = true;
  $: usernameValid = usernameRegex.test(username);
  let error: APIError | null = null;

  const changeUsername = async () => {
    try {
      const resp = await apiFetchClient<MeUser>("/users/@me", "PATCH", { username });

      data.user = resp;
      userStore.set(resp);
      localStorage.setItem("pronouns-user", JSON.stringify(resp));
      error = null;
    } catch (e) {
      error = e as APIError;
    }
  };

  let deleteOpen = false;
  const toggleDeleteOpen = () => (deleteOpen = !deleteOpen);
  let deleteUsername = "";
  let deleteError: APIError | null = null;

  const deleteAccount = async () => {
    try {
      await fastFetchClient("/users/@me", "DELETE");

      userStore.set(null);
      localStorage.removeItem("pronouns-token");
      localStorage.removeItem("pronouns-user");
      toggleDeleteOpen();
      addToast({ header: "Deleted account", body: "Your account is now pending deletion." });
      goto("/");
    } catch (e) {
      deleteUsername = "";
      deleteError = e as APIError;
    }
  };

  let invalidateModalOpen = false;
  const toggleInvalidateModalOpen = () => (invalidateModalOpen = !invalidateModalOpen);
  let invalidateError: APIError | null = null;

  const invalidateAllTokens = async () => {
    try {
      await fastFetchClient("/auth/tokens", "DELETE");

      invalidateError = null;
      userStore.set(null);
      localStorage.removeItem("pronouns-token");
      localStorage.removeItem("pronouns-user");
      toggleInvalidateModalOpen();
      addToast({
        header: "Invalidated tokens",
        body: "Invalidated all your tokens, please log in again.",
      });
      goto("/");
    } catch (e) {
      invalidateError = e as APIError;
    }
  };

  const DONATE_ALERT_STORE = "alert-1681976313";
  let donateAlertOpen = false;
  const closeDonateAlert = () => {
    donateAlertOpen = false;
    localStorage.setItem(DONATE_ALERT_STORE, "dismissed");
  };

  onMount(() => {
    if (!localStorage.getItem(DONATE_ALERT_STORE)) {
      donateAlertOpen = true;
    }
  });
</script>

<Alert color="secondary" fade={false} isOpen={donateAlertOpen} toggle={closeDonateAlert}>
  If you find pronouns.cc useful and have the means, I would really appreciate a <a
    href="https://liberapay.com/u1f320/"
    target="_blank">donation</a
  >
  to keep the site running, fast, and ad-free!
  <br />
  It's not required and doesn't give you any perks, but running this website takes time and money so
  it would really help.
  <span class="text-muted">
    (If you want to hide this alert, press the <Icon name="x-lg" alt="X" /> in the top right to hide
    it on this computer or phone)
  </span>
</Alert>

<h1>Your profile</h1>

<div class="grid">
  <div class="row">
    <div class="col-lg">
      <h4>Username</h4>
      <div class="input-group m-1 w-75">
        <input type="text" class="form-control" bind:value={username} />
        <Button
          color="secondary"
          on:click={() => changeUsername()}
          disabled={username === data.user.name || !usernameValid}>Change username</Button
        >
      </div>
      {#if username !== data.user.name}
        <p class="text-muted">
          <Icon name="info-circle-fill" aria-hidden /> Changing your username will make any existing
          links to your or your members' profiles invalid.
          <br />
          Your username must be unique, be at most 40 characters long, and only contain letters from
          the basic English alphabet, dashes, underscores, and periods. Your username is used as part
          of your profile link, you can set a separate display name.
          {#if !usernameValid}
            <br />
            <span class="text-danger-emphasis"
              ><Icon name="exclamation-triangle-fill" aria-hidden /> That username is not valid.</span
            >
          {/if}
        </p>
      {/if}
      {#if error}
        <ErrorAlert {error} />
      {/if}
    </div>
    <div class="col-lg-4">
      <FallbackImage width={200} urls={userAvatars(data.user)} alt="Your avatar" />
      <p>
        To change your avatar, go to <a href="/edit/profile">edit profile</a>.
      </p>
    </div>
  </div>
  <div class="row">
    <div class="col">
      <h4>Force log out</h4>
      <p>
        If you think one of your tokens might have been compromised, you can log out on all devices
        by clicking this button.
      </p>
      <p>
        <Button color="danger" on:click={toggleInvalidateModalOpen}>Force log out</Button>
      </p>
      <Modal isOpen={invalidateModalOpen} toggle={toggleInvalidateModalOpen}>
        <ModalHeader toggle={toggleInvalidateModalOpen}>Force log out</ModalHeader>
        <ModalBody>
          <p>If you want to force log out on all devices, click the button below.</p>
          {#if invalidateError}
            <ErrorAlert error={invalidateError} />
          {/if}
        </ModalBody>
        <ModalFooter>
          <Button color="danger" on:click={invalidateAllTokens}>Force log out</Button>
          <Button color="secondary" on:click={toggleInvalidateModalOpen}>Cancel</Button>
        </ModalFooter>
      </Modal>
    </div>
  </div>
  <div class="row">
    <div class="col">
      <h4>Account info</h4>
      <Table bordered striped hover>
        <tbody>
          <tr>
            <th scope="row">ID</th>
            <td><code>{data.user.id}</code></td>
          </tr>
          <tr>
            <th scope="row">Members</th>
            <td>{data.user.members.length}/{MAX_MEMBERS}</td>
          </tr>
          {#if data.invitesEnabled}
            <tr>
              <th scope="row">Invites</th>
              <td>{data.invites.length}/{data.user.max_invites}</td>
            </tr>
          {/if}
        </tbody>
      </Table>
    </div>
  </div>
  <div class="row">
    <div class="col">
      <h4>Delete account</h4>
      <p>If you want to initiate the account deletion process, click the button below:</p>
      <p>
        <Button color="danger" on:click={toggleDeleteOpen}>Delete account</Button>
      </p>
      <p class="text-muted">
        <Icon name="info-circle-fill" /> Your account will be deleted 30 days after initiating the deletion
        process. If you want to cancel your account deletion within 30 days, log in again and confirm
        that you want to cancel deletion.
        <strong
          >Simply logging in again will not cancel account deletion. Also, your account cannot be
          recovered after the deletion period has passed.</strong
        >
      </p>
      <Modal isOpen={deleteOpen} toggle={toggleDeleteOpen}>
        <ModalHeader toggle={toggleDeleteOpen}>Delete account</ModalHeader>
        <ModalBody>
          <p>If you want to delete your account, type your current username below:</p>
          <p>
            <input type="text" class="form-control" bind:value={deleteUsername} />
          </p>
          {#if deleteError}
            <ErrorAlert error={deleteError} />
          {/if}
        </ModalBody>
        <ModalFooter>
          <Button
            color="danger"
            disabled={deleteUsername !== data.user.name}
            on:click={deleteAccount}
          >
            Delete account
          </Button>
          <Button color="secondary" on:click={toggleDeleteOpen}>Cancel</Button>
        </ModalFooter>
      </Modal>
    </div>
  </div>
</div>
