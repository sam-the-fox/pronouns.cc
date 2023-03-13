<script lang="ts">
  import { onMount } from "svelte";
  import { Alert, Icon } from "sveltestrap";

  import { goto } from "$app/navigation";
  import type { APIError, MeUser } from "$lib/api/entities";
  import { apiFetch } from "$lib/api/fetch";
  import { userStore } from "$lib/store";
  import type { PageData } from "./$types";

  interface SignupResponse {
    user: MeUser;
    token: string;
  }

  export let data: PageData;

  onMount(() => {
    if (data.token && data.user) {
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
      const resp = await apiFetch<SignupResponse>("/auth/discord/signup", {
        method: "POST",
        body: {
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
</script>

<svelte:head>
  <title>Log in with Discord - pronouns.cc</title>
</svelte:head>

<h1>Log in with Discord</h1>

{#if data.error}
  <Alert color="danger" fade={false}>
    <h4 class="alert-heading">An error occurred</h4>
    <b>{data.error.code}:</b>
    {data.error.message}
  </Alert>
{/if}
{#if data.ticket}
  <form on:submit|preventDefault={signupForm}>
    <div>
      <label for="discord">Discord username</label>
      <input id="discord" class="form-control" name="discord" disabled value={data.discord} />
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
          <Icon name="info" /> You currently need an invite code to sign up. You can get one from an
          existing user.
        </div>
      </div>
    {/if}
    <div class="form-text">
      By signing up, you agree to the <a href="/page/tos">terms of service</a> and the
      <a href="/page/privacy">privacy policy</a>.
    </div>
    <button type="submit" class="btn btn-primary">Sign up</button>
  </form>
{:else}
  Loading...
{/if}
