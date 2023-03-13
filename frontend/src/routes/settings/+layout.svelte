<script lang="ts">
  import { page } from "$app/stores";
  import type { LayoutData } from "./$types";

  import { ListGroup, ListGroupItem } from "sveltestrap";

  export let data: LayoutData;
</script>

<div class="grid">
  <div class="row">
    <div class="col-md-3 m-3">
      <h1>Settings</h1>

      <ListGroup>
        <ListGroupItem tag="a" active={$page.url.pathname === "/settings"} href="/settings">
          Your profile
        </ListGroupItem>
        {#if data.require_invite}
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
          API tokens
        </ListGroupItem>
        <ListGroupItem
          tag="a"
          active={$page.url.pathname === "/settings/export"}
          href="/settings/export"
        >
          Data export
        </ListGroupItem>
      </ListGroup>
    </div>
    <div class="col-md m-3">
      <slot />
    </div>
  </div>
</div>
