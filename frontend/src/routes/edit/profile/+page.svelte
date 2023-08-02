<script lang="ts">
  import {
    MAX_DESCRIPTION_LENGTH,
    userAvatars,
    type APIError,
    type Field,
    type FieldEntry,
    type MeUser,
    type Pronoun,
    PreferenceSize,
    type CustomPreferences,
    type PrideFlag,
  } from "$lib/api/entities";
  import FallbackImage from "$lib/components/FallbackImage.svelte";
  import { userStore } from "$lib/store";
  import {
    Alert,
    Button,
    ButtonGroup,
    Card,
    CardBody,
    CardHeader,
    FormGroup,
    InputGroup,
    Icon,
    Input,
    Popover,
    TabContent,
    TabPane,
    InputGroupText,
  } from "sveltestrap";
  import { encode } from "base64-arraybuffer";
  import { DateTime, FixedOffsetZone } from "luxon";
  import { apiFetchClient } from "$lib/api/fetch";
  import { PUBLIC_SHORT_BASE } from "$env/static/public";
  import IconButton from "$lib/components/IconButton.svelte";
  import EditableField from "../EditableField.svelte";
  import EditableName from "../EditableName.svelte";
  import EditablePronouns from "../EditablePronouns.svelte";
  import ErrorAlert from "$lib/components/ErrorAlert.svelte";
  import { addToast, delToast } from "$lib/toast";
  import type { PageData, Snapshot } from "./$types";
  import { charCount, renderMarkdown } from "$lib/utils";
  import MarkdownHelp from "../MarkdownHelp.svelte";
  import prettyBytes from "pretty-bytes";
  import CustomPreference from "./CustomPreference.svelte";
  import FlagButton from "../FlagButton.svelte";

  const MAX_AVATAR_BYTES = 1_000_000;

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
  let avatar_files: FileList | null;

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
  $: getAvatar(avatar_files).then((b64) => (avatar = b64));

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

  const getAvatar = async (list: FileList | null) => {
    if (!list || list.length === 0) return null;
    if (list[0].size > MAX_AVATAR_BYTES) {
      addToast({
        header: "Avatar too large",
        body: `This avatar is too large, please resize it (maximum is ${prettyBytes(
          MAX_AVATAR_BYTES,
        )}, the file you tried to upload is ${prettyBytes(list[0].size)})`,
      });
      return null;
    }

    const buffer = await list[0].arrayBuffer();
    const base64 = encode(buffer);

    const uri = `data:${list[0].type};base64,${base64}`;
    console.log(uri.slice(0, 128));

    return uri;
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

  const moveName = (index: number, up: boolean) => {
    if (up && index == 0) return;
    if (!up && index == names.length - 1) return;

    const newIndex = up ? index - 1 : index + 1;

    const temp = names[index];
    names[index] = names[newIndex];
    names[newIndex] = temp;
  };

  const movePronoun = (index: number, up: boolean) => {
    if (up && index == 0) return;
    if (!up && index == pronouns.length - 1) return;

    const newIndex = up ? index - 1 : index + 1;

    const temp = pronouns[index];
    pronouns[index] = pronouns[newIndex];
    pronouns[newIndex] = temp;
  };

  const moveField = (index: number, up: boolean) => {
    if (up && index == 0) return;
    if (!up && index == fields.length - 1) return;

    const newIndex = up ? index - 1 : index + 1;

    const temp = fields[index];
    fields[index] = fields[newIndex];
    fields[newIndex] = temp;
  };

  const moveLink = (index: number, up: boolean) => {
    if (up && index == 0) return;
    if (!up && index == links.length - 1) return;

    const newIndex = up ? index - 1 : index + 1;

    const temp = links[index];
    links[index] = links[newIndex];
    links[newIndex] = temp;
  };

  const moveFlag = (index: number, up: boolean) => {
    if (up && index == 0) return;
    if (!up && index == flags.length - 1) return;

    const newIndex = up ? index - 1 : index + 1;

    const temp = flags[index];
    flags[index] = flags[newIndex];
    flags[newIndex] = temp;
  };

  const addFlag = (flag: PrideFlag) => {
    flags = [...flags, flag];
  };

  const removeFlag = (index: number) => {
    flags.splice(index, 1);
    flags = [...flags];
  };

  const addName = (event: Event) => {
    event.preventDefault();

    names = [...names, { value: newName, status: "okay" }];
    newName = "";
  };

  const addPronouns = (event: Event) => {
    event.preventDefault();

    if (newPronouns in data.pronouns) {
      const fullSet = data.pronouns[newPronouns];
      pronouns = [
        ...pronouns,
        {
          pronouns: fullSet.pronouns.join("/"),
          display_text: fullSet.display || null,
          status: "okay",
        },
      ];
    } else {
      pronouns = [...pronouns, { pronouns: newPronouns, display_text: null, status: "okay" }];
    }
    newPronouns = "";
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

  const removeName = (index: number) => {
    names.splice(index, 1);
    names = [...names];
  };

  const removePronoun = (index: number) => {
    pronouns.splice(index, 1);
    pronouns = [...pronouns];
  };

  const removeLink = (index: number) => {
    links.splice(index, 1);
    links = [...links];
  };

  const removeField = (index: number) => {
    fields.splice(index, 1);
    fields = [...fields];
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
  <TabPane tabId="avatar" tab="Names and avatar" active>
    <div class="row mt-3">
      <div class="col-md">
        <div class="row">
          <div class="col-md text-center">
            {#if avatar === ""}
              <FallbackImage alt="Current avatar" urls={[]} width={200} />
            {:else if avatar}
              <img
                width={200}
                height={200}
                src={avatar}
                alt="New avatar"
                class="rounded-circle img-fluid"
              />
            {:else}
              <FallbackImage alt="Current avatar" urls={userAvatars(data.user)} width={200} />
            {/if}
          </div>
          <div class="col-md">
            <input
              class="form-control"
              id="avatar"
              type="file"
              bind:files={avatar_files}
              accept="image/png, image/jpeg, image/gif, image/webp"
            />
            <p class="text-muted mt-3">
              <Icon name="info-circle-fill" aria-hidden /> Only PNG, JPEG, GIF, and WebP images can be
              used as avatars. Avatars cannot be larger than 1 MB, and animated avatars will be made
              static.
            </p>
            <p>
              <!-- svelte-ignore a11y-invalid-attribute -->
              <a href="" on:click={() => (avatar = "")}>Remove avatar</a>
            </p>
          </div>
        </div>
      </div>
      <div class="col-md">
        <FormGroup floating label="Username">
          <Input bind:value={data.user.name} readonly />
          <p class="text-muted mt-1">
            <Icon name="info-circle-fill" aria-hidden />
            You can change your username in
            <a href="/settings" class="text-reset">your settings</a>.
          </p>
        </FormGroup>
        <FormGroup floating label="Display name">
          <Input bind:value={display_name} />
          <p class="text-muted mt-1">
            <Icon name="info-circle-fill" aria-hidden />
            Your display name is used in page titles and as a header.
          </p>
        </FormGroup>
      </div>
    </div>
    <div>
      <h4>Names</h4>
      {#each names as _, index}
        <EditableName
          bind:value={names[index].value}
          bind:status={names[index].status}
          preferences={custom_preferences}
          moveUp={() => moveName(index, true)}
          moveDown={() => moveName(index, false)}
          remove={() => removeName(index)}
        />
      {/each}
      <form class="input-group m-1" on:submit={addName}>
        <input type="text" class="form-control" bind:value={newName} />
        <IconButton type="submit" color="success" icon="plus" tooltip="Add name" />
      </form>
    </div>
  </TabPane>
  <TabPane tabId="bio" tab="Bio">
    <div class="mt-3">
      <div class="form">
        <textarea class="form-control" style="height: 200px;" bind:value={bio} />
      </div>
      <p class="text-muted mt-1">
        Using {charCount(bio)}/{MAX_DESCRIPTION_LENGTH} characters
      </p>
      <p class="text-muted my-2">
        <MarkdownHelp />
      </p>
      {#if bio}
        <hr />
        <Card>
          <CardHeader>Preview</CardHeader>
          <CardBody>
            {@html renderMarkdown(bio)}
          </CardBody>
        </Card>
      {/if}
    </div>
  </TabPane>
  <TabPane tabId="pronouns" tab="Pronouns">
    <div class="mt-3">
      <div class="col-md">
        {#each pronouns as _, index}
          <EditablePronouns
            bind:pronoun={pronouns[index]}
            preferences={custom_preferences}
            moveUp={() => movePronoun(index, true)}
            moveDown={() => movePronoun(index, false)}
            remove={() => removePronoun(index)}
          />
        {/each}
        <form class="input-group m-1" on:submit={addPronouns}>
          <input
            type="text"
            class="form-control"
            placeholder="New pronouns"
            bind:value={newPronouns}
            required
          />
          <IconButton
            type="submit"
            color="success"
            icon="plus"
            tooltip="Add pronouns"
            disabled={newPronouns === ""}
          />
          <Button id="pronouns-help" color="secondary"><Icon name="question" /></Button>
          <Popover target="pronouns-help" placement="bottom">
            For common pronouns, the short form (e.g. "she/her" or "he/him") is enough; for less
            common pronouns, you will have to use all five forms (e.g. "ce/cir/cir/cirs/cirself").
          </Popover>
        </form>
      </div>
    </div>
  </TabPane>
  <TabPane tabId="fields" tab="Fields">
    {#if data.user.fields.length === 0}
      <Alert class="mt-3" color="secondary" fade={false}>
        Fields are extra categories you can add separate from names and pronouns.<br />
        For example, you could use them for gender terms, honorifics, or compliments.
      </Alert>
    {/if}
    <div class="grid gap-3">
      <div class="row row-cols-1 row-cols-md-2">
        {#each fields as _, index}
          <EditableField
            bind:field={fields[index]}
            preferences={custom_preferences}
            deleteField={() => removeField(index)}
            moveField={(up) => moveField(index, up)}
          />
        {/each}
      </div>
    </div>
    <div>
      <Button on:click={() => (fields = [...fields, { name: "New field", entries: [] }])}>
        <Icon name="plus" aria-hidden /> Add new field
      </Button>
    </div>
  </TabPane>
  <TabPane tabId="flags" tab="Flags">
    <div class="mt-3">
      {#each flags as _, index}
        <ButtonGroup class="m-1">
          <IconButton
            icon="chevron-left"
            color="secondary"
            tooltip="Move flag to the left"
            click={() => moveFlag(index, true)}
          />
          <IconButton
            icon="chevron-right"
            color="secondary"
            tooltip="Move flag to the right"
            click={() => moveFlag(index, false)}
          />
          <FlagButton
            flag={flags[index]}
            tooltip="Remove this flag from your profile"
            on:click={() => removeFlag(index)}
          />
        </ButtonGroup>
      {/each}
    </div>
    <hr />
    <div class="row">
      <div class="col-md">
        <Input
          placeholder="Filter flags"
          bind:value={flagSearch}
          disabled={data.flags.length === 0}
        />
        <div class="p-2">
          {#each filteredFlags as flag (flag.id)}
            <FlagButton
              {flag}
              tooltip="Add this flag to your profile"
              on:click={() => addFlag(flag)}
            />
          {:else}
            {#if data.flags.length === 0}
              You haven't uploaded any flags yet.
            {:else}
              There are no flags matching your search <strong>{flagSearch}</strong>.
            {/if}
          {/each}
        </div>
      </div>
      <div class="col-md">
        <Alert color="secondary" fade={false}>
          {#if data.flags.length === 0}
            <p><strong>Why can't I see any flags?</strong></p>
            <p>
              There are thousands of pride flags, and it would be impossible to bundle all of them
              by default. Many labels also have multiple different flags that are favoured by
              different people. Because of this, there are no flags available by default--instead,
              you can upload flags in your <a href="/settings/flags">settings</a>. Your main profile
              and your member profiles can all have different flags.
            </p>
          {:else}
            To upload and delete flags, go to your <a href="/settings/flags">settings</a>.
          {/if}
        </Alert>
      </div>
    </div>
  </TabPane>
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
            <InputGroupText><Icon name="clock" aria-hidden /></InputGroupText>
            <Input disabled value={timezone !== null ? timezone : "Unset"} />
            <Button on:click={detectTimezone}>Detect</Button>
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
