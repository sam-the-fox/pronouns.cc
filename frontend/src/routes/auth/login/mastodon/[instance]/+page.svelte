<script lang="ts">
  import { onMount } from "svelte";
  import { Alert, Button, Icon } from "sveltestrap";

  import { goto } from "$app/navigation";
  import type { APIError, MeUser } from "$lib/api/entities";
  import { apiFetch } from "$lib/api/fetch";
  import { userStore } from "$lib/store";
  import type { PageData } from "./$types";
  import ErrorAlert from "$lib/components/ErrorAlert.svelte";

  interface SignupResponse {
    user: MeUser;
    token: string;
  }

  export let data: PageData;

  onMount(() => {
    if (!data.is_deleted && data.token && data.user) {
      localStorage.setItem("pronouns-token", data.token);
      localStorage.setItem("pronouns-user", JSON.stringify(data.user));
      userStore.set(data.user);
      goto("/");
    }
  });

  let username = "";
  let invite = "";

  const signupForm = async () => {
    try {
      const resp = await apiFetch<SignupResponse>("/auth/mastodon/signup", {
        method: "POST",
        body: {
          instance: data.instance,
          ticket: data.ticket,
          username: username,
          invite_code: invite,
        },
      });

      localStorage.setItem("pronouns-token", resp.token);
      localStorage.setItem("pronouns-user", JSON.stringify(resp.user));
      userStore.set(resp.user);
      goto("/");
    } catch (e) {
      data.error = e as APIError;
    }
  };

  let deleteCancelled = false;
  let deleteError: APIError | null = null;
  const cancelDelete = async () => {
    try {
      await apiFetch<any>("/auth/cancel-delete", {
        method: "GET",
        headers: {
          "X-Delete-Token": data.token!,
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
  <title>Log in with Mastodon - pronouns.cc</title>
</svelte:head>

<h1>Log in with Mastodon</h1>

{#if data.error}
  <ErrorAlert error={data.error} />
{/if}
{#if data.ticket}
  <form on:submit|preventDefault={signupForm}>
    <div>
      <label for="fediverse">Fediverse username</label>
      <input
        id="fediverse"
        class="form-control"
        name="fediverse"
        disabled
        value="{data.fediverse}@{data.instance}"
      />
    </div>
    <div>
      <label for="username">Username</label>
      <input id="username" class="form-control" name="username" bind:value={username} />
    </div>
    {#if data.require_invite}
      <div>
        <label for="invite">Invite code</label>
        <input
          id="invite"
          class="form-control"
          name="invite"
          bind:value={invite}
          aria-describedby="invite-help"
        />
        <div id="invite-help" class="form-text">
          <Icon name="info-circle-fill" /> You currently need an invite code to sign up. You can get
          one from an existing user.
        </div>
      </div>
    {/if}
    <div class="form-text">
      By signing up, you agree to the <a href="/page/tos">terms of service</a> and the
      <a href="/page/privacy">privacy policy</a>.
    </div>
    <Button type="submit" color="primary">Sign up</Button>
  </form>
{:else if data.is_deleted && data.token}
  <p>Your account is pending deletion since {data.deleted_at}.</p>
  <p>If you wish to cancel deletion, press the button below.</p>
  <p>
    <Button color="primary" on:click={cancelDelete} disabled={deleteCancelled}
      >Cancel account deletion</Button
    >
  </p>
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
