<script lang="ts">
  import { getContext } from "svelte";
  import type { Writable } from "svelte/store";
  import { MAX_DESCRIPTION_LENGTH, type MeUser } from "$lib/api/entities";
  import { charCount, renderMarkdown } from "$lib/utils";
  import MarkdownHelp from "$lib/components/edit/MarkdownHelp.svelte";
  import { Card, CardBody, CardHeader } from "sveltestrap";

  const user = getContext<Writable<MeUser>>("user");
</script>

<div class="form">
  <textarea class="form-control" style="height: 200px;" bind:value={$user.bio} />
</div>
<p class="text-muted mt-1">
  Using {charCount($user.bio || "")}/{MAX_DESCRIPTION_LENGTH} characters
</p>
<p class="text-muted my-2">
  <MarkdownHelp />
</p>
{#if $user.bio}
  <hr />
  <Card>
    <CardHeader>Preview</CardHeader>
    <CardBody>
      {@html renderMarkdown($user.bio)}
    </CardBody>
  </Card>
{/if}
