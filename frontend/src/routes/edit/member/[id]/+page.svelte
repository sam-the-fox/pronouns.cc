<script lang="ts">
  import { goto } from "$app/navigation";
  import {
    MAX_DESCRIPTION_LENGTH,
    memberAvatars,
    type APIError,
    type Field,
    type FieldEntry,
    type Member,
    type Pronoun,
    type PrideFlag,
  } from "$lib/api/entities";
  import FallbackImage from "$lib/components/FallbackImage.svelte";
  import {
    Button,
    ButtonGroup,
    FormGroup,
    Icon,
    Input,
    Modal,
    ModalBody,
    ModalFooter,
    Popover,
    TabContent,
    TabPane,
    Card,
    CardBody,
    CardHeader,
    Alert,
  } from "sveltestrap";
  import { encode } from "base64-arraybuffer";
  import prettyBytes from "pretty-bytes";
  import { apiFetchClient, fastFetchClient } from "$lib/api/fetch";
  import IconButton from "$lib/components/IconButton.svelte";
  import EditableField from "../../EditableField.svelte";
  import EditableName from "../../EditableName.svelte";
  import EditablePronouns from "../../EditablePronouns.svelte";
  import ErrorAlert from "$lib/components/ErrorAlert.svelte";
  import type { PageData } from "./$types";
  import { addToast, delToast } from "$lib/toast";
  import { memberNameRegex } from "$lib/api/regex";
  import { charCount, renderMarkdown } from "$lib/utils";
  import MarkdownHelp from "../../MarkdownHelp.svelte";
  import FlagButton from "../../FlagButton.svelte";

  const MAX_AVATAR_BYTES = 1_000_000;

  export let data: PageData;

  if (data.user.id !== data.member.user.id) {
    addToast({ header: "Not your member", body: "You cannot edit another user's member." });
    goto(`/@${data.member.user.name}/${data.member.name}`);
  }

  let error: APIError | null = null;

  let bio: string = data.member.bio || "";
  let name: string = data.member.name;
  let display_name: string = data.member.display_name || "";
  let links: string[] = window.structuredClone(data.member.links);
  let names: FieldEntry[] = window.structuredClone(data.member.names);
  let pronouns: Pronoun[] = window.structuredClone(data.member.pronouns);
  let fields: Field[] = window.structuredClone(data.member.fields);
  let flags: PrideFlag[] = window.structuredClone(data.member.flags);
  let unlisted: boolean = data.member.unlisted || false;

  let memberNameValid = true;
  $: memberNameValid = memberNameRegex.test(name);

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

  let modified = false;

  $: modified = isModified(
    data.member,
    bio,
    name,
    display_name,
    links,
    names,
    pronouns,
    fields,
    flags,
    avatar,
    unlisted,
  );
  $: getAvatar(avatar_files).then((b64) => (avatar = b64));

  const isModified = (
    member: Member,
    bio: string,
    name: string,
    display_name: string,
    links: string[],
    names: FieldEntry[],
    pronouns: Pronoun[],
    fields: Field[],
    flags: PrideFlag[],
    avatar: string | null,
    unlisted: boolean,
  ) => {
    if (name !== member.name) return true;
    if (bio !== member.bio) return true;
    if (display_name !== member.display_name) return true;
    if (!linksEqual(links, member.links)) return true;
    if (!fieldsEqual(fields, member.fields)) return true;
    if (!flagsEqual(flags, member.flags)) return true;
    if (!namesEqual(names, member.names)) return true;
    if (!pronounsEqual(pronouns, member.pronouns)) return true;
    if (avatar !== null) return true;
    if (unlisted !== member.unlisted) return true;

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

    return uri;
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

  const updateMember = async () => {
    const toastId = addToast({
      header: "Saving changes",
      body: "Saving changes, please wait...",
      duration: -1,
    });

    try {
      const resp = await apiFetchClient<Member>(`/members/${data.member.id}`, "PATCH", {
        name,
        display_name,
        avatar,
        bio,
        links,
        names,
        pronouns,
        fields,
        flags: flags.map((flag) => flag.id),
        unlisted,
      });

      addToast({ header: "Success", body: "Successfully saved changes!" });

      data.member = resp;
      avatar = null;
      error = null;
    } catch (e) {
      error = e as APIError;
    } finally {
      delToast(toastId);
    }
  };

  const deleteMember = async () => {
    try {
      await fastFetchClient(`/members/${data.member.id}`, "DELETE");

      toggleDeleteOpen();
      addToast({
        header: "Deleted member",
        body: `Successfully deleted member ${data.member.name}!`,
      });
      goto(`/@${data.member.user.name}`);
    } catch (e) {
      deleteName = "";
      deleteError = e as APIError;
    }
  };

  let deleteModalPronoun = "the member's";
  $: deleteModalPronoun = updateModalPronoun(pronouns);

  const updateModalPronoun = (pronouns: Pronoun[]) => {
    const filtered = pronouns.filter((entry) => entry.status === "favourite");
    if (filtered.length < 1) return "the member's";

    const split = filtered[0].pronouns.split("/");
    if (split.length !== 5) return "the member's";
    return split[2];
  };

  let deleteOpen = false;
  const toggleDeleteOpen = () => (deleteOpen = !deleteOpen);
  let deleteName = "";
  let deleteError: APIError | null = null;
</script>

<svelte:head>
  <title>Edit member profile - pronouns.cc</title>
</svelte:head>

<h1>
  Edit member profile
  <ButtonGroup>
    <IconButton
      color="secondary"
      icon="chevron-left"
      href="/@{data.member.user.name}/{data.member.name}"
      tooltip="Back to member"
    />
    {#if modified}
      <Button color="success" on:click={() => updateMember()} disabled={!memberNameValid}
        >Save changes</Button
      >
    {/if}
    <Button color="danger" on:click={toggleDeleteOpen}
      >Delete {data.member.display_name ?? data.member.name}</Button
    >
  </ButtonGroup>
</h1>

<Modal header="Delete member" isOpen={deleteOpen} toggle={toggleDeleteOpen}>
  <ModalBody>
    <p>
      If you want to delete this member, type {deleteModalPronoun} name (<code
        >{data.member.name}</code
      >) below:
    </p>
    <p>
      <input type="text" class="form-control" bind:value={deleteName} />
    </p>
    {#if deleteError}
      <ErrorAlert error={deleteError} />
    {/if}
  </ModalBody>
  <ModalFooter>
    <Button color="danger" disabled={deleteName !== data.member.name} on:click={deleteMember}>
      Delete member
    </Button>
    <Button color="secondary" on:click={toggleDeleteOpen}>Cancel</Button>
  </ModalFooter>
</Modal>

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
              <FallbackImage alt="Current avatar" urls={memberAvatars(data.member)} width={200} />
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
        <FormGroup floating label="Name">
          <Input bind:value={name} />
          <p class="text-muted mt-1">
            <Icon name="info-circle-fill" aria-hidden />
            The member name is only used as part of the link to their profile page.
          </p>
        </FormGroup>
        {#if !memberNameValid}
          <p class="text-danger-emphasis mb-2">That member name is not valid.</p>
        {/if}
        <FormGroup floating label="Display name">
          <Input bind:value={display_name} />
        </FormGroup>
        <p class="text-muted mt-1">
          <Icon name="info-circle-fill" aria-hidden />
          Your display name is used in page titles and as a header.
        </p>
      </div>
    </div>
    <div>
      <h4>Names</h4>
      {#each names as _, index}
        <EditableName
          bind:value={names[index].value}
          bind:status={names[index].status}
          preferences={data.user.custom_preferences}
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
            preferences={data.user.custom_preferences}
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
    {#if data.member.fields.length === 0}
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
            preferences={data.user.custom_preferences}
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
  <TabPane tabId="other" tab="Other">
    <div class="row mt-3">
      <div class="col-md">
        <div class="form-check">
          <input class="form-check-input" type="checkbox" bind:checked={unlisted} id="unlisted" />
          <label class="form-check-label" for="unlisted">Hide from member list</label>
        </div>
        <p class="text-muted mt-1">
          {#if data.user.list_private}
            <Icon name="exclamation-triangle-fill" aria-hidden />
            Your member list is currently hidden, so <strong>this setting has no effect</strong>. If
            you want to make your member list visible again,
            <a href="/edit/profile">edit your user profile</a>.
            <br />
          {/if}
          <Icon name="info-circle-fill" aria-hidden />
          This <em>only</em> hides this member from your member list.
          <strong>
            This member will still be visible to anyone at
            <code class="text-nowrap">pronouns.cc/@{data.user.name}/{data.member.name}</code>.
          </strong>
        </p>
      </div>
    </div>
  </TabPane>
</TabContent>
