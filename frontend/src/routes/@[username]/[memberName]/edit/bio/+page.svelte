<script lang="ts">
  import { getContext } from "svelte";
  import type { Writable } from "svelte/store";
  import { MAX_DESCRIPTION_LENGTH, type Member } from "$lib/api/entities";
  import { charCount, renderMarkdown } from "$lib/utils";
  import MarkdownHelp from "$lib/components/edit/MarkdownHelp.svelte";
  import { Card, CardBody, CardHeader } from "@sveltestrap/sveltestrap";

  const member = getContext<Writable<Member>>("member");
</script>

<div class="form">
  <textarea class="form-control" style="height: 200px;" bind:value={$member.bio} />
</div>
<p class="text-muted mt-1">
  Using {charCount($member.bio || "")}/{MAX_DESCRIPTION_LENGTH} characters
</p>
<p class="text-muted my-2">
  <MarkdownHelp />
</p>
{#if $member.bio}
  <hr />
  <Card>
    <CardHeader>Preview</CardHeader>
    <CardBody>
      {@html renderMarkdown($member.bio)}
    </CardBody>
  </Card>
{/if}
