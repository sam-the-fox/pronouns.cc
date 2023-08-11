<script lang="ts">
  import { page } from "$app/stores";
    import { ErrorCode } from "$lib/api/entities";
</script>

<h1>An error occurred ({$page.status})</h1>

{#if $page.error?.code === ErrorCode.NotFound}
  <p>
    The page you were looking for was not found. If you're sure the page exists, check for any typos
    in the address.
  </p>
{:else if $page.error?.code === ErrorCode.Forbidden || $page.error?.code === ErrorCode.InvalidToken}
  <p>You're not logged in, or you aren't allowed to access this page.</p>
{:else if $page.error?.code === 429}
  <p>You've exceeded a rate limit, please try again later.</p>
{:else if $page.error?.code === 500}
  <p>An internal error occurred. Please try again later.</p>
  <p>
    If this error keeps happening, please <a
      href="https://codeberg.org/pronounscc/pronouns.cc/issues"
      target="_blank"
      rel="noreferrer">file a bug report</a
    > with an explanation of what you did to cause the error.
  </p>
{/if}

<p>Error code: <code>{$page.error?.code}</code></p>
