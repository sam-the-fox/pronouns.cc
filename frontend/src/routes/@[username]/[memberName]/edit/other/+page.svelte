<script lang="ts">
  import { getContext } from "svelte";
  import type { Writable } from "svelte/store";
  import { DateTime } from "luxon";
  import { Button, ButtonGroup, Icon } from "@sveltestrap/sveltestrap";

  import type { APIError, Member } from "$lib/api/entities";
  import { PUBLIC_SHORT_BASE } from "$env/static/public";
  import IconButton from "$lib/components/IconButton.svelte";
  import { apiFetchClient } from "$lib/api/fetch";
  import { addToast } from "$lib/toast";
  import type { PageData } from "./$types";
  import ErrorAlert from "$lib/components/ErrorAlert.svelte";

  export let data: PageData;

  const member = getContext<Writable<Member>>("member");

  let error: APIError | null = null;

  const now = DateTime.now().toLocal();
  let canRerollSid: boolean;
  $: canRerollSid =
    now.diff(DateTime.fromISO(data.user.last_sid_reroll).toLocal(), "hours").hours >= 1;

  const rerollSid = async () => {
    try {
      const resp = await apiFetchClient<Member>(`/members/${data.member.id_new}/reroll`);
      addToast({ header: "Success", body: "Rerolled short ID!" });
      error = null;
      $member.sid = resp.sid;
    } catch (e) {
      error = e as APIError;
    }
  };

  const copyShortURL = async () => {
    const url = `${PUBLIC_SHORT_BASE}/${data.member.sid}`;
    await navigator.clipboard.writeText(url);
    addToast({ body: "Copied the short link to your clipboard!", duration: 2000 });
  };
</script>

{#if error}
  <ErrorAlert {error} />
{/if}

<div class="row">
  <div class="col-md">
    <div class="form-check">
      <input
        class="form-check-input"
        type="checkbox"
        bind:checked={$member.unlisted}
        id="unlisted"
      />
      <label class="form-check-label" for="unlisted">Hide from member list</label>
    </div>
    <p class="text-muted mt-1">
      {#if data.user.list_private}
        <Icon name="exclamation-triangle-fill" aria-hidden />
        Your member list is currently hidden, so <strong>this setting has no effect</strong>. If you
        want to make your member list visible again,
        <a href="/@{$member.user.name}/other">edit your user profile</a>.
        <br />
      {/if}
      <Icon name="info-circle-fill" aria-hidden />
      This <em>only</em> hides this member from your member list.
      <strong>
        This member will still be visible to anyone at
        <code class="text-nowrap">pronouns.cc/@{$member.user.name}/{$member.name}</code>.
      </strong>
    </p>
  </div>
  {#if PUBLIC_SHORT_BASE}
    <div class="col-md">
      <p>
        Current short ID: <code>{$member.sid}</code>
        <ButtonGroup class="mb-1">
          <Button color="secondary" disabled={!canRerollSid} on:click={() => rerollSid()}
            >Reroll short ID</Button
          >
          <IconButton
            icon="link-45deg"
            tooltip="Copy short link"
            color="secondary"
            click={copyShortURL}
          />
        </ButtonGroup>
        <br />
        <span class="text-muted">
          <Icon name="info-circle-fill" aria-hidden />
          This ID is used in <code>prns.cc</code> links. You can reroll one short ID every hour (shared
          between your main profile and all members) by pressing the button above.
        </span>
      </p>
    </div>
  {/if}
</div>
