<script lang="ts">
  import type { PageData } from "./$types";

  export let data: PageData;
</script>

<svelte:head>
  <title>Login - pronouns.cc</title>
</svelte:head>

<div class="container">
  <h1>Log in or sign up</h1>
  <div class="row">
    <div class="col">
      <p>
        <a class="btn btn-primary" href={data.discord} role="button">Log in with Discord</a>
      </p>
    </div>
    <div class="col">

    </div>
  </div>
</div>
