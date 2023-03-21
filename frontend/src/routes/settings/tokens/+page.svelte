<script lang="ts">
  import { DateTime } from "luxon";
  import { Icon, Table } from "sveltestrap";
  import type { PageData } from "./$types";

  export let data: PageData;

  import { decodeJwt } from "jose";
  const claims = decodeJwt(localStorage.getItem("pronouns-token")!);
</script>

<h1>Tokens ({data.tokens.length})</h1>

<Table bordered striped hover>
  <thead>
    <th>ID</th>
    <th>Created at</th>
    <th>Expires at</th>
    <th>Current?</th>
  </thead>
  <tbody>
    {#each data.tokens as token}
      <tr>
        <td><code>{token.id}</code></td>
        <td>{DateTime.fromISO(token.created).toLocal().toLocaleString(DateTime.DATETIME_MED)}</td>
        <td>{DateTime.fromISO(token.expires).toLocal().toLocaleString(DateTime.DATETIME_MED)}</td>
        <td
          >{#if claims["jti"] === token.id}<Icon name="check-lg" alt="Current token" />{:else}<Icon
              name="x-lg"
              alt="Not current token"
            />{/if}</td
        >
      </tr>
    {/each}
  </tbody>
</Table>
