<script lang="ts">
  import { getContext } from "svelte";
  import type { Writable } from "svelte/store";

  import { PreferenceSize, type APIError, type MeUser } from "$lib/api/entities";
  import IconButton from "$lib/components/IconButton.svelte";
  import { Button, ButtonGroup, FormGroup, Icon, Input, InputGroup } from "sveltestrap";
  import { PUBLIC_SHORT_BASE } from "$env/static/public";
  import CustomPreference from "./CustomPreference.svelte";
  import { DateTime, FixedOffsetZone } from "luxon";
  import { addToast } from "$lib/toast";
  import ErrorAlert from "$lib/components/ErrorAlert.svelte";
  import { apiFetchClient } from "$lib/api/fetch";

  const user = getContext<Writable<MeUser>>("user");
  let error: APIError | null = null;

  // Custom preferences code
  let preferenceIds: string[];
  $: preferenceIds = Object.keys($user.custom_preferences);

  const addPreference = () => {
    const id = crypto.randomUUID();

    $user.custom_preferences[id] = {
      icon: "question",
      tooltip: "New preference",
      size: PreferenceSize.Normal,
      muted: false,
      favourite: false,
    };
    $user.custom_preferences = $user.custom_preferences;
  };

  const removePreference = (id: string) => {
    delete $user.custom_preferences[id];
    $user.custom_preferences = $user.custom_preferences;
  };

  // Timezone code
  let currentTime = "";
  let displayTimezone = "";
  $: setTime($user.timezone);

  const detectTimezone = () => {
    $user.timezone = DateTime.local().zoneName;
  };

  const setTime = (timezone: string | null) => {
    if (!timezone) {
      currentTime = "";
      displayTimezone = "";
      return;
    }

    const offset = DateTime.now().setZone(timezone).offset;
    const zone = FixedOffsetZone.instance(offset);

    currentTime = now.setZone(zone).toLocaleString(DateTime.TIME_SIMPLE);
    displayTimezone = zone.formatOffset(now.toUnixInteger(), "narrow");
  };

  // SID code
  const now = DateTime.now().toLocal();
  let canRerollSid: boolean;
  $: canRerollSid = now.diff(DateTime.fromISO($user.last_sid_reroll).toLocal(), "hours").hours >= 1;

  const copyShortURL = async () => {
    const url = `${PUBLIC_SHORT_BASE}/${$user.sid}`;
    await navigator.clipboard.writeText(url);
    addToast({ body: "Copied the short link to your clipboard!", duration: 2000 });
  };

  const rerollSid = async () => {
    try {
      const resp = await apiFetchClient<MeUser>("/users/@me/reroll");
      addToast({ header: "Success", body: "Rerolled short ID!" });
      error = null;
      $user.sid = resp.sid;
    } catch (e) {
      error = e as APIError;
    }
  };
</script>

{#if error}
  <ErrorAlert {error} />
{/if}

<div class="row">
  <div class="col-md">
    <FormGroup floating label={'"Members" header text'}>
      <Input bind:value={$user.member_title} placeholder="Members" />
      <p class="text-muted mt-1">
        <Icon name="info-circle-fill" aria-hidden />
        This is the text used for the "Members" heading. If you leave it blank, the default text will
        be used.
      </p>
    </FormGroup>
    {#if PUBLIC_SHORT_BASE}
      <hr />
      <p>
        Current short ID: <code>{$user.sid}</code>
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
    {/if}
  </div>
  <div class="col-md">
    <div class="form-check">
      <input
        class="form-check-input"
        type="checkbox"
        bind:checked={$user.list_private}
        id="listPrivate"
      />
      <label class="form-check-label" for="listPrivate">Hide member list</label>
    </div>
    <p class="text-muted mt-1">
      <Icon name="info-circle-fill" aria-hidden />
      This only hides your member <em>list</em>.
      <strong>
        Your members will still be visible to anyone at
        <code class="text-nowrap">pronouns.cc/@{$user.name}/[member-name]</code>.
      </strong>
    </p>
    <hr />
    <div class="m-1">
      <p class="mt-1 my-2">
        You can optionally set your timezone, which will show your current local time on your
        profile.
      </p>
      <InputGroup>
        <Button on:click={detectTimezone}>Detect timezone</Button>
        <Input disabled value={$user.timezone !== null ? $user.timezone : "Unset"} />
        <Button on:click={() => ($user.timezone = null)}>Reset</Button>
      </InputGroup>
      <p class="mt-2">
        {#if $user.timezone}
          This will show up on your profile like this:
          <Icon name="clock" aria-hidden />
          {currentTime} <span class="text-body-secondary">(UTC{displayTimezone})</span>
          <br />
        {/if}
        <span class="text-muted">
          Your timezone is never shared directly, only the difference between UTC and your current
          timezone is.
        </span>
      </p>
    </div>
  </div>
</div>
<div>
  <h3>
    Preferences <Button on:click={addPreference} color="success"
      ><Icon name="plus" aria-hidden /> Add new</Button
    >
  </h3>
  {#each preferenceIds as id}
    <CustomPreference
      bind:preference={$user.custom_preferences[id]}
      remove={() => removePreference(id)}
    />
  {/each}
</div>
