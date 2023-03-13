<script>
  import { ListGroup, ListGroupItem } from "sveltestrap";
</script>

<div class="grid">
  <div class="row">
    <div class="col-md-3 m-3">
      <h1>Settings</h1>

      <ListGroup>
        <ListGroupItem tag="a" href="/settings">Your profile</ListGroupItem>
        <ListGroupItem tag="a" href="/settings/invites">Invites</ListGroupItem>
        <ListGroupItem tag="a" href="/settings/tokens">API tokens</ListGroupItem>
      </ListGroup>
    </div>
    <div class="col-md m-3">
      <slot />
    </div>
  </div>
</div>
