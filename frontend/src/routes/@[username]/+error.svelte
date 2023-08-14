<script lang="ts">
  import { page } from "$app/stores";
  import { ErrorCode } from "$lib/api/entities";
</script>

<svelte:head>
  <title>Error - pronouns.cc</title>
</svelte:head>

{#if $page.error?.code === ErrorCode.Forbidden || $page.error?.code === ErrorCode.InvalidToken}
  <h1>Not logged in</h1>
  <p>
    Either you aren't logged in, or your login has expired. Please <a href="/auth/login"
      >log in again</a
    >.
  </p>
{:else if $page.error?.code === ErrorCode.NotOwnMember}
  <h1>Not your member</h1>
  <p>You can only edit your own members.</p>
{:else}
  <h1>An error occurred ({$page.status})</h1>

  {#if $page.status === 404}
    <p>The user you were looking for couldn't be found. Please check for any typos.</p>
  {:else if $page.status === 429}
    <p>You've exceeded a rate limit, please try again later.</p>
  {:else if $page.status === 500}
    <p>An internal error occurred. Please try again later.</p>
    <p>
      If this error keeps happening, please <a
        href="https://codeberg.org/pronounscc/pronouns.cc/issues"
        target="_blank"
        rel="noreferrer">file a bug report</a
      > with an explanation of what you did to cause the error.
    </p>
  {/if}

  <p>Error message: <code>{$page.error?.message}</code></p>
{/if}
