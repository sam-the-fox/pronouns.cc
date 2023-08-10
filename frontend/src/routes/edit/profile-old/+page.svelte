<script lang="ts">
  import {
    type APIError,
    type Field,
    type FieldEntry,
    type MeUser,
    type Pronoun,
    PreferenceSize,
    type CustomPreferences,
    type PrideFlag,
  } from "$lib/api/entities";
  import { userStore } from "$lib/store";
  import {
    Button,
    ButtonGroup,
    FormGroup,
    InputGroup,
    Icon,
    Input,
    TabContent,
    TabPane,
  } from "sveltestrap";
  import { DateTime, FixedOffsetZone } from "luxon";
  import { apiFetchClient } from "$lib/api/fetch";
  import { PUBLIC_SHORT_BASE } from "$env/static/public";
  import IconButton from "$lib/components/IconButton.svelte";
  import ErrorAlert from "$lib/components/ErrorAlert.svelte";
  import { addToast, delToast } from "$lib/toast";
  import type { PageData, Snapshot } from "./$types";
  import CustomPreference from "./CustomPreference.svelte";

  export let data: PageData;

  let error: APIError | null = null;

  let bio: string = data.user.bio || "";
  let display_name: string = data.user.display_name || "";
  let member_title: string = data.user.member_title || "";
  let links: string[] = window.structuredClone(data.user.links);
  let names: FieldEntry[] = window.structuredClone(data.user.names);
  let pronouns: Pronoun[] = window.structuredClone(data.user.pronouns);
  let fields: Field[] = window.structuredClone(data.user.fields);
  let flags: PrideFlag[] = window.structuredClone(data.user.flags);
  let list_private = data.user.list_private;
  let custom_preferences = window.structuredClone(data.user.custom_preferences);
  let timezone = data.user.timezone;

  let avatar: string | null;

  let newName = "";
  let newPronouns = "";
  let newLink = "";

  let flagSearch = "";
  let filteredFlags: PrideFlag[];
  $: filteredFlags = filterFlags(flagSearch, data.flags);

  const filterFlags = (search: string, flags: PrideFlag[]) => {
    return (
      search
        ? flags.filter((flag) => flag.name.toLocaleLowerCase().includes(search.toLocaleLowerCase()))
        : flags
    ).slice(0, 25);
  };

  let preferenceIds: string[];
  $: preferenceIds = Object.keys(custom_preferences);

  let modified = false;

  $: modified = isModified(
    data.user,
    bio,
    display_name,
    links,
    names,
    pronouns,
    fields,
    flags,
    avatar,
    member_title,
    list_private,
    custom_preferences,
    timezone,
  );

  const isModified = (
    user: MeUser,
    bio: string,
    display_name: string,
    links: string[],
    names: FieldEntry[],
    pronouns: Pronoun[],
    fields: Field[],
    flags: PrideFlag[],
    avatar: string | null,
    member_title: string,
    list_private: boolean,
    custom_preferences: CustomPreferences,
    timezone: string | null,
  ) => {
    if (bio !== (user.bio || "")) return true;
    if (display_name !== (user.display_name || "")) return true;
    if (member_title !== (user.member_title || "")) return true;
    if (!linksEqual(links, user.links)) return true;
    if (!fieldsEqual(fields, user.fields)) return true;
    if (!flagsEqual(flags, user.flags)) return true;
    if (!namesEqual(names, user.names)) return true;
    if (!pronounsEqual(pronouns, user.pronouns)) return true;
    if (!customPreferencesEqual(custom_preferences, user.custom_preferences)) return true;
    if (avatar !== null) return true;
    if (list_private !== user.list_private) return true;
    if (timezone !== user.timezone) return true;

    return false;
  };

  const fieldsEqual = (arr1: Field[], arr2: Field[]) => {
    if (arr1?.length !== arr2?.length) return false;
    if (!arr1.every((_, i) => arr1[i].entries.length === arr2[i].entries.length)) return false;
    if (!arr1.every((_, i) => arr1[i].name === arr2[i].name)) return false;

    return arr1.every((_, i) =>
      arr1[i].entries.every(
        (entry, j) =>
          entry.value === arr2[i].entries[j].value && entry.status === arr2[i].entries[j].status,
      ),
    );
  };

  const namesEqual = (arr1: FieldEntry[], arr2: FieldEntry[]) => {
    if (arr1?.length !== arr2?.length) return false;
    if (!arr1.every((_, i) => arr1[i].value === arr2[i].value)) return false;
    if (!arr1.every((_, i) => arr1[i].status === arr2[i].status)) return false;

    return true;
  };

  const pronounsEqual = (arr1: Pronoun[], arr2: Pronoun[]) => {
    if (arr1?.length !== arr2?.length) return false;
    if (!arr1.every((_, i) => arr1[i].pronouns === arr2[i].pronouns)) return false;
    if (!arr1.every((_, i) => arr1[i].display_text === arr2[i].display_text)) return false;
    if (!arr1.every((_, i) => arr1[i].status === arr2[i].status)) return false;

    return true;
  };

  const linksEqual = (arr1: string[], arr2: string[]) => {
    if (arr1.length !== arr2.length) return false;
    return arr1.every((_, i) => arr1[i] === arr2[i]);
  };

  const flagsEqual = (arr1: PrideFlag[], arr2: PrideFlag[]) => {
    if (arr1.length !== arr2.length) return false;
    return arr1.every((_, i) => arr1[i].id === arr2[i].id);
  };

  const customPreferencesEqual = (obj1: CustomPreferences, obj2: CustomPreferences) => {
    if (Object.keys(obj2).some((key) => !(key in obj1))) return false;

    return Object.keys(obj1)
      .map((key) => {
        if (!(key in obj2)) return false;
        return (
          obj1[key].icon === obj2[key].icon &&
          obj1[key].tooltip === obj2[key].tooltip &&
          obj1[key].favourite === obj2[key].favourite &&
          obj1[key].muted === obj2[key].muted &&
          obj1[key].size === obj2[key].size
        );
      })
      .every((entry) => entry);
  };

  let currentTime = "";
  let displayTimezone = "";
  $: setTime(timezone);

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

  const moveLink = (index: number, up: boolean) => {
    if (up && index == 0) return;
    if (!up && index == links.length - 1) return;

    const newIndex = up ? index - 1 : index + 1;

    const temp = links[index];
    links[index] = links[newIndex];
    links[newIndex] = temp;
  };

  const addLink = (event: Event) => {
    event.preventDefault();

    links = [...links, newLink];
    newLink = "";
  };

  const addPreference = () => {
    const id = crypto.randomUUID();

    custom_preferences[id] = {
      icon: "question",
      tooltip: "New preference",
      size: PreferenceSize.Normal,
      muted: false,
      favourite: false,
    };
    custom_preferences = custom_preferences;
  };

  const removeLink = (index: number) => {
    links.splice(index, 1);
    links = [...links];
  };

  const removePreference = (id: string) => {
    delete custom_preferences[id];
    custom_preferences = custom_preferences;
  };

  const updateUser = async () => {
    const toastId = addToast({
      header: "Saving changes",
      body: "Saving changes, please wait...",
      duration: -1,
    });

    try {
      const resp = await apiFetchClient<MeUser>("/users/@me", "PATCH", {
        display_name,
        avatar,
        bio,
        links,
        names,
        pronouns,
        fields,
        member_title,
        list_private,
        timezone: timezone || "",
        custom_preferences,
        flags: flags.map((flag) => flag.id),
      });

      data.user = resp;
      custom_preferences = resp.custom_preferences;
      userStore.set(resp);
      localStorage.setItem("pronouns-user", JSON.stringify(resp));

      addToast({ header: "Success", body: "Successfully saved changes!" });

      avatar = null;
      error = null;
    } catch (e) {
      error = e as APIError;
    } finally {
      delToast(toastId);
    }
  };

  const now = DateTime.now().toLocal();
  let canRerollSid: boolean;
  $: canRerollSid =
    now.diff(DateTime.fromISO(data.user.last_sid_reroll).toLocal(), "hours").hours >= 1;

  const rerollSid = async () => {
    try {
      const resp = await apiFetchClient<MeUser>("/users/@me/reroll");
      addToast({ header: "Success", body: "Rerolled short ID!" });
      error = null;
      data.user.sid = resp.sid;
    } catch (e) {
      error = e as APIError;
    }
  };

  const copyShortURL = async () => {
    const url = `${PUBLIC_SHORT_BASE}/${data.user.sid}`;
    await navigator.clipboard.writeText(url);
    addToast({ body: "Copied the short link to your clipboard!", duration: 2000 });
  };

  const detectTimezone = () => {
    timezone = DateTime.local().zoneName;
  };

  interface SnapshotData {
    bio: string;
    display_name: string;
    member_title: string;
    links: string[];
    names: FieldEntry[];
    pronouns: Pronoun[];
    fields: Field[];
    flags: PrideFlag[];
    list_private: boolean;
    timezone: string | null;
    custom_preferences: CustomPreferences;

    avatar: string | null;
    newName: string;
    newPronouns: string;
    newLink: string;
  }

  export const snapshot: Snapshot<SnapshotData> = {
    capture: () => ({
      bio,
      display_name,
      member_title,
      links,
      names,
      pronouns,
      fields,
      flags,
      list_private,
      timezone,
      custom_preferences,
      avatar,
      newName,
      newPronouns,
      newLink,
    }),
    restore: (value) => {
      bio = value.bio;
      display_name = value.display_name;
      member_title = value.member_title;
      links = value.links;
      names = value.names;
      pronouns = value.pronouns;
      fields = value.fields;
      flags = value.flags;
      list_private = value.list_private;
      timezone = value.timezone;
      custom_preferences = value.custom_preferences;
      avatar = value.avatar;
      newName = value.newName;
      newPronouns = value.newPronouns;
      newLink = value.newLink;
    },
  };
</script>

<svelte:head>
  <title>Edit profile - pronouns.cc</title>
</svelte:head>

<h1>
  Edit profile
  <ButtonGroup>
    <Button color="secondary" href="/@{data.user.name}">
      <Icon name="chevron-left" />
      Back to your profile
    </Button>
    {#if modified}
      <Button color="success" on:click={() => updateUser()}>Save changes</Button>
    {/if}
  </ButtonGroup>
</h1>

{#if error}
  <ErrorAlert {error} />
{/if}

<TabContent>
  <TabPane tabId="links" tab="Links">
    <div class="mt-3">
      {#each links as _, index}
        <div class="input-group m-1">
          <IconButton
            icon="chevron-up"
            color="secondary"
            tooltip="Move link up"
            click={() => moveLink(index, true)}
          />
          <IconButton
            icon="chevron-down"
            color="secondary"
            tooltip="Move link down"
            click={() => moveLink(index, false)}
          />
          <input type="text" class="form-control" bind:value={links[index]} />
          <IconButton
            color="danger"
            icon="trash3"
            tooltip="Remove link"
            click={() => removeLink(index)}
          />
        </div>
      {/each}
      <form class="input-group m-1" on:submit={addLink}>
        <input type="text" class="form-control" bind:value={newLink} />
        <IconButton type="submit" color="success" icon="plus" tooltip="Add link" />
      </form>
    </div>
  </TabPane>
  <TabPane tabId="other" tab="Preferences & other">
    <div class="row mt-3">
      <div class="col-md">
        <FormGroup floating label={'"Members" header text'}>
          <Input bind:value={member_title} placeholder="Members" />
          <p class="text-muted mt-1">
            <Icon name="info-circle-fill" aria-hidden />
            This is the text used for the "Members" heading. If you leave it blank, the default text
            will be used.
          </p>
        </FormGroup>
        {#if PUBLIC_SHORT_BASE}
          <hr />
          <p>
            Current short ID: <code>{data.user.sid}</code>
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
            bind:checked={list_private}
            id="listPrivate"
          />
          <label class="form-check-label" for="listPrivate">Hide member list</label>
        </div>
        <p class="text-muted mt-1">
          <Icon name="info-circle-fill" aria-hidden />
          This only hides your member <em>list</em>.
          <strong>
            Your members will still be visible to anyone at
            <code class="text-nowrap">pronouns.cc/@{data.user.name}/[member-name]</code>.
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
            <Input disabled value={timezone !== null ? timezone : "Unset"} />
            <Button on:click={() => (timezone = null)}>Reset</Button>
          </InputGroup>
          <p class="mt-2">
            {#if timezone}
              This will show up on your profile like this:
              <Icon name="clock" aria-hidden />
              {currentTime} <span class="text-body-secondary">(UTC{displayTimezone})</span>
              <br />
            {/if}
            <span class="text-muted">
              Your timezone is never shared directly, only the difference between UTC and your
              current timezone is.
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
          bind:preference={custom_preferences[id]}
          remove={() => removePreference(id)}
        />
      {/each}
    </div>
  </TabPane>
</TabContent>
