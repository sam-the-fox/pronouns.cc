<script lang="ts">
  import { page } from "$app/stores";
  import type { LayoutData } from "./$types";
  import { Button, ListGroup, ListGroupItem, Modal, ModalBody, ModalFooter } from "sveltestrap";
  import { userStore } from "$lib/store";
  import { goto } from "$app/navigation";
  import { addToast } from "$lib/toast";

  export let data: LayoutData;

  let isOpen = false;
  const toggle = () => (isOpen = !isOpen);

  const logOut = () => {
    userStore.set(null);
    localStorage.removeItem("pronouns-token");
    localStorage.removeItem("pronouns-user");
    toggle();

    addToast({ header: "Logged out", body: "Successfully logged out!" });
    goto("/");
  };
</script>

<svelte:head>
  <title>Settings - pronouns.cc</title>
</svelte:head>

<div class="grid">
  <div class="row">
    <div class="col-md-3 m-3">
      <h1>Settings</h1>

      <ListGroup>
        <ListGroupItem tag="a" active={$page.url.pathname === "/settings"} href="/settings">
          Your profile
        </ListGroupItem>
        <ListGroupItem
          tag="a"
          active={$page.url.pathname === "/settings/auth"}
          href="/settings/auth"
        >
          Authentication
        </ListGroupItem>
        {#if data.invitesEnabled}
          <ListGroupItem
            tag="a"
            active={$page.url.pathname === "/settings/invites"}
            href="/settings/invites"
          >
            Invites
          </ListGroupItem>
        {/if}
        <ListGroupItem
          tag="a"
          active={$page.url.pathname === "/settings/tokens"}
          href="/settings/tokens"
        >
          Tokens
        </ListGroupItem>
        <ListGroupItem
          tag="a"
          active={$page.url.pathname === "/settings/export"}
          href="/settings/export"
        >
          Data export
        </ListGroupItem>
        <ListGroupItem tag="button" on:click={toggle}>Log out</ListGroupItem>
      </ListGroup>
    </div>
    <div class="col-md m-3">
      <slot />
    </div>
  </div>
  <Modal header="Log out" {isOpen} {toggle}>
    <ModalBody>Are you sure you want to log out?</ModalBody>
    <ModalFooter>
      <Button color="primary" on:click={logOut}>Log out</Button>
      <Button color="secondary" on:click={toggle}>Cancel</Button>
    </ModalFooter>
  </Modal>
</div>
