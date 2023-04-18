<script lang="ts">
  import { goto } from "$app/navigation";
  import type { APIError } from "$lib/api/entities";
  import { apiFetch } from "$lib/api/fetch";
  import ErrorAlert from "$lib/components/ErrorAlert.svelte";
  import { userStore } from "$lib/store";
  import { addToast } from "$lib/toast";
  import { onMount } from "svelte";
  import {
    Button,
    Icon,
    Input,
    ListGroup,
    ListGroupItem,
    Modal,
    ModalBody,
    ModalFooter,
  } from "sveltestrap";
  import type { PageData } from "./$types";
  import fediverse from "./fediverse.svg";

  export let data: PageData;

  let error: APIError | null = null;
  let instance = "";
  let fediDisabled = false;

  let modalOpen = false;
  let toggleModal = () => (modalOpen = !modalOpen);

  onMount(() => {
    if ($userStore) {
      addToast({ header: "Error", body: "You are already logged in." });
      goto("/");
    }
  });

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

<svelte:head>
  <title>Login - pronouns.cc</title>
</svelte:head>

<div class="container">
  <h1>Log in or sign up</h1>
  <div class="row">
    <div class="col-md-4 mb-1">
      <ListGroup>
        <ListGroupItem tag="button" on:click={toggleModal}>Log in with the Fediverse</ListGroupItem>
        <ListGroupItem tag="a" href={data.discord}>Log in with Discord</ListGroupItem>
        <ListGroupItem tag="a" href={data.tumblr}>Log in with Tumblr</ListGroupItem>
      </ListGroup>
      <Modal header="Pick an instance" isOpen={modalOpen} toggle={toggleModal}>
        <ModalBody>
          <Input placeholder="Instance (e.g. mastodon.social)" bind:value={instance} />
          <p class="text-muted mt-2">
            <Icon name="info-circle-fill" aria-label="Info" /> This should be the domain you use to access
            posts. For example, if your username is <code>@timmie@mastodon.example</code>, but your
            user <em>page</em> is at
            <code>social.mastodon.example/timmie</code>, you should fill in
            <code>social.mastodon.example</code>.
          </p>
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
      <p class="text-muted mt-2 mx-1">
        <Icon name="info-circle-fill" aria-hidden /> By signing in, you consent to pronouns.cc storing
        a token in your browser to identify your account.
      </p>
    </div>
    <div class="col-md">
      <p>
        <b>Choose an authentication provider to get started.</b> You can add more providers later.
      </p>
    </div>
  </div>
</div>
