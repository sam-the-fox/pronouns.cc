<script lang="ts">
  import { goto } from "$app/navigation";
  import type { APIError, MeUser } from "$lib/api/entities";
  import { fastFetch } from "$lib/api/fetch";
  import { usernameRegex } from "$lib/api/regex";
  import ErrorAlert from "$lib/components/ErrorAlert.svelte";
  import { PUBLIC_HCAPTCHA_SITEKEY } from "$env/static/public";
  import HCaptcha from "svelte-hcaptcha";
  import { userStore } from "$lib/store";
  import { addToast } from "$lib/toast";
  import { DateTime } from "luxon";
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
  } from "@sveltestrap/sveltestrap";

  export let authType: string;
  export let remoteName: string | undefined;
  export let error: APIError | undefined;
  export let requireInvite: boolean | undefined;
  export let requireCaptcha: boolean | undefined;
  export let isDeleted: boolean | undefined;
  export let ticket: string | undefined;
  export let token: string | undefined;
  export let user: MeUser | undefined;
  export let deletedAt: string | undefined;
  export let selfDelete: boolean | undefined;
  export let deleteReason: string | undefined;

  onMount(() => {
    if (!isDeleted && token && user) {
      localStorage.setItem("pronouns-token", token);
      localStorage.setItem("pronouns-user", JSON.stringify(user));
      userStore.set(user);
      addToast({ header: "Logged in", body: "Successfully logged in!" });
      goto(`/@${user.name}`);
    }
  });

  let deleteCancelled: boolean;
  let deleteError: APIError | null;

  let username: string;
  let usernameValid = true;
  $: usernameValid = usernameRegex.test(username);

  let inviteCode: string;
  let forceDeleteName = "";
  let forceDeleteModalOpen = false;
  let toggleForceDeleteModal = () => (forceDeleteModalOpen = !forceDeleteModalOpen);

  export let linkAccount: () => Promise<void>;
  export let signupForm: (
    username: string,
    inviteCode: string,
    captchaToken: string,
  ) => Promise<void>;

  let captchaToken = "";
  let captcha: any;

  const captchaSuccess = (token: any) => {
    captchaToken = token.detail.token;
  };

  let canSubmit = false;
  $: canSubmit = usernameValid && (!!captchaToken || !requireCaptcha);

  const captchaError = () => {
    addToast({
      header: "Captcha failed",
      body: "There was an error verifying the captcha, please try again.",
    });
    captcha.reset();
  };

  export const resetCaptcha = (): void => captcha.reset();

  const forceDeleteAccount = async () => {
    try {
      await fastFetch("/auth/force-delete", {
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
      await fastFetch("/auth/cancel-delete", {
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
  <title>Log in with {authType} - pronouns.cc</title>
</svelte:head>

<h1>Log in with {authType}</h1>

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
  <form on:submit|preventDefault={() => signupForm(username, inviteCode, captchaToken)}>
    <div>
      <FormGroup floating label="{authType} username">
        <Input readonly value={remoteName} />
      </FormGroup>
    </div>
    <div>
      <FormGroup floating label="Username">
        <Input id="username" name="username" bind:value={username} />
      </FormGroup>
      <div id="username-help" class="text-muted">
        <Icon name="info-circle-fill" aria-label="Info" /> Your username must be unique, be at most 40
        characters long, and only contain letters from the basic English alphabet, dashes, underscores,
        and periods. Your username is used as part of your profile link, you can set a separate display
        name.
      </div>
    </div>
    {#if requireInvite}
      <div>
        <FormGroup floating label="Invite code">
          <Input id="invite" name="invite" aria-describedby="invite-help" bind:value={inviteCode} />
        </FormGroup>
        <div id="invite-help" class="text-muted">
          <Icon name="info-circle-fill" aria-label="Info" /> You currently need an invite code to sign
          up. You can get one from an existing user.
        </div>
      </div>
    {/if}
    {#if requireCaptcha}
      <div class="mt-2 mx-2 mb-1">
        <HCaptcha
          bind:this={captcha}
          sitekey={PUBLIC_HCAPTCHA_SITEKEY}
          on:success={captchaSuccess}
          on:error={captchaError}
        />
      </div>
    {/if}
    <div class="text-muted my-1">
      By signing up, you agree to the <a href="/page/terms">terms of service</a> and the
      <a href="/page/privacy">privacy policy</a>.
    </div>
    <p>
      <Button type="submit" color="primary" disabled={!canSubmit}>Sign up</Button>
      {#if !usernameValid && username.length > 0}
        <span class="text-danger-emphasis mb-2">That username is not valid.</span>
      {/if}
    </p>
  </form>
{:else if isDeleted && token && selfDelete && deletedAt}
  <p>
    Your account is pending deletion since {DateTime.fromISO(deletedAt)
      .toLocal()
      .toLocaleString(DateTime.DATETIME_MED)}.
  </p>
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
  {#if deleteCancelled}
    <Alert color="secondary" fade={false}>
      Account deletion cancelled! You can now <a href="/auth/login">log in</a> again.
    </Alert>
  {/if}
  {#if deleteError}
    <ErrorAlert error={deleteError} />
  {/if}
{:else if isDeleted && token && !selfDelete && deletedAt}
  <p>
    Your account is pending deletion since {DateTime.fromISO(deletedAt)
      .toLocal()
      .toLocaleString(DateTime.DATETIME_MED)}.
  </p>
  <p>
    <strong>Your account was deactivated by a moderator.</strong> You cannot cancel deletion. The moderator
    gave the following reason:
  </p>
  <blockquote class="blockquote">
    {deleteReason}
  </blockquote>
  <p>
    Your account will be fully deleted 180 days after being deactivated. If you want your data wiped
    immediately instead, press the force delete link below.
  </p>
  <p>
    <Button color="link" on:click={toggleForceDeleteModal}>Force delete account</Button>
  </p>
{:else}
  Loading...
{/if}

<Modal header="Force delete account" isOpen={forceDeleteModalOpen} toggle={toggleForceDeleteModal}>
  <ModalBody>
    <p>
      If you want to delete your account, type your username (<code>{user?.name}</code>) below:
      <br />
      <b>
        This is irreversible! Your account <i>cannot</i> be recovered after you press "Force delete account".
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
