<script lang="ts">
  import { goto } from "$app/navigation";
  import {
    MAX_DESCRIPTION_LENGTH,
    memberAvatars,
    WordStatus,
    type APIError,
    type Field,
    type FieldEntry,
    type Member,
    type Pronoun,
  } from "$lib/api/entities";
  import FallbackImage from "$lib/components/FallbackImage.svelte";
  import { userStore } from "$lib/store";
  import {
    Alert,
    Button,
    ButtonGroup,
    FormGroup,
    Icon,
    Input,
    Modal,
    ModalBody,
    ModalFooter,
  } from "sveltestrap";
  import { encode } from "base64-arraybuffer";
  import { apiFetchClient } from "$lib/api/fetch";
  import IconButton from "$lib/components/IconButton.svelte";
  import EditableField from "../../EditableField.svelte";
  import EditableName from "../../EditableName.svelte";
  import EditablePronouns from "../../EditablePronouns.svelte";
  import ErrorAlert from "$lib/components/ErrorAlert.svelte";
  import type { PageData } from "./$types";
  import { addToast, delToast } from "$lib/toast";

  const MAX_AVATAR_BYTES = 1_000_000;

  export let data: PageData;

  if (!$userStore || $userStore.id !== data.member.user.id) {
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

  let avatar: string | null;
  let avatar_files: FileList | null;

  let newName = "";
  let newPronouns = "";
  let newPronounsDisplay = "";
  let newLink = "";

  let modified = false;

  $: modified = isModified(bio, name, display_name, links, names, pronouns, fields, avatar);
  $: getAvatar(avatar_files).then((b64) => (avatar = b64));

  const isModified = (
    bio: string,
    name: string,
    display_name: string,
    links: string[],
    names: FieldEntry[],
    pronouns: Pronoun[],
    fields: Field[],
    avatar: string | null,
  ) => {
    if (name !== data.member.name) return true;
    if (bio !== data.member.bio) return true;
    if (display_name !== data.member.display_name) return true;
    if (!linksEqual(links, data.member.links)) return true;
    if (!fieldsEqual(fields, data.member.fields)) return true;
    if (!namesEqual(names, data.member.names)) return true;
    if (!pronounsEqual(pronouns, data.member.pronouns)) return true;
    if (avatar !== null) return true;

    return false;
  };

  const fieldsEqual = (arr1: Field[], arr2: Field[]) => {
    if (arr1?.length !== arr2?.length) return false;
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

  const getAvatar = async (list: FileList | null) => {
    if (!list || list.length === 0) return null;
    if (list[0].size > MAX_AVATAR_BYTES) return null;

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

  const addName = () => {
    names = [...names, { value: newName, status: WordStatus.Okay }];
    newName = "";
  };

  const addPronouns = () => {
    pronouns = [
      ...pronouns,
      { pronouns: newPronouns, display_text: newPronounsDisplay || null, status: WordStatus.Okay },
    ];
    newPronouns = "";
    newPronounsDisplay = "";
  };

  const addLink = () => {
    links = [...links, newLink];
    newLink = "";
  };

  const removeName = (index: number) => {
    if (names.length === 1) names = [];
    else if (index === 0) names = names.slice(1);
    else if (index === names.length - 1) names = names.slice(0, names.length - 1);
    else names = [...names.slice(0, index - 1), ...names.slice(0, index + 1)];
  };

  const removePronoun = (index: number) => {
    if (pronouns.length === 1) pronouns = [];
    else if (index === 0) pronouns = pronouns.slice(1);
    else if (index === pronouns.length - 1) pronouns = pronouns.slice(0, pronouns.length - 1);
    else pronouns = [...pronouns.slice(0, index - 1), ...pronouns.slice(0, index + 1)];
  };

  const removeLink = (index: number) => {
    if (links.length === 1) links = [];
    else if (index === 0) links = links.slice(1);
    else if (index === links.length - 1) links = links.slice(0, links.length - 1);
    else links = [...links.slice(0, index - 1), ...links.slice(0, index + 1)];
  };

  const removeField = (index: number) => {
    if (fields.length === 1) fields = [];
    else if (index === 0) fields = fields.slice(1);
    else if (index === fields.length - 1) fields = fields.slice(0, fields.length - 1);
    else fields = [...fields.slice(0, index - 1), ...fields.slice(0, index + 1)];
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
      });

      addToast({ header: "Success", body: "Successfully saved changes!" });

      data.member = resp;
      avatar = null;
      error = null;
      modified = false;
    } catch (e) {
      error = e as APIError;
    } finally {
      delToast(toastId);
    }
  };

  const deleteMember = async () => {
    try {
      await apiFetchClient<any>(`/members/${data.member.id}`, "DELETE");

      toggleDeleteOpen();
      goto(`/@${data.member.user.name}`);
    } catch (e) {
      deleteName = "";
      deleteError = e as APIError;
    }
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
      <Button color="success" on:click={() => updateMember()}>Save changes</Button>
    {/if}
    <Button color="danger" on:click={toggleDeleteOpen}
      >Delete {data.member.display_name ?? data.member.name}</Button
    >
  </ButtonGroup>
</h1>

<Modal header="Delete member" isOpen={deleteOpen} toggle={toggleDeleteOpen}>
  <ModalBody>
    <p>
      If you want to delete this member, type their name (<code>{data.member.name}</code>) below:
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

{#if !$userStore}
  <Alert color="danger" fade={false}>Error: No user object</Alert>
{:else}
  <div class="grid">
    <div class="row m-1">
      <div class="col-md">
        <h4>Avatar</h4>
        <div class="row">
          <div class="col-md text-center">
            {#if avatar}
              <img width={200} src={avatar} alt="New avatar" class="rounded-circle img-fluid" />
            {:else}
              <FallbackImage alt="Current avatar" urls={memberAvatars(data.member)} width={200} />
            {/if}
          </div>
          <div class="col-md mt-2">
            <input
              class="form-control"
              id="avatar"
              type="file"
              bind:files={avatar_files}
              accept="image/png, image/jpeg, image/gif, image/webp"
            />
            <p class="text-muted mt-3">
              <Icon name="info-circle-fill" /> Only PNG, JPEG, GIF, and WebP can be used as avatars.
              Avatars cannot be larger than 1 MB, and animated avatars will be made static.
            </p>
          </div>
        </div>
      </div>
      <div class="col-md">
        <div>
          <FormGroup floating label="Name">
            <Input bind:value={name} />
          </FormGroup>
        </div>
        <div>
          <FormGroup floating label="Display name">
            <Input bind:value={display_name} />
          </FormGroup>
        </div>
        <div>
          <FormGroup floating label="Bio ({bio.length}/{MAX_DESCRIPTION_LENGTH})">
            <textarea style="min-height: 100px;" class="form-control" bind:value={bio} />
          </FormGroup>
        </div>
      </div>
    </div>
    <div class="row m-1">
      <div class="col-md">
        <h4>Names</h4>
        {#each names as _, index}
          <EditableName
            bind:value={names[index].value}
            bind:status={names[index].status}
            moveUp={() => moveName(index, true)}
            moveDown={() => moveName(index, false)}
            remove={() => removeName(index)}
          />
        {/each}
        <div class="input-group m-1">
          <input type="text" class="form-control" bind:value={newName} />
          <IconButton color="success" icon="plus" tooltip="Add name" click={() => addName()} />
        </div>
      </div>
      <div class="col-md">
        <h4>Links</h4>
        {#each links as _, index}
          <div class="input-group m-1">
            <input type="text" class="form-control" bind:value={links[index]} />
            <IconButton
              color="danger"
              icon="trash3"
              tooltip="Remove link"
              click={() => removeLink(index)}
            />
          </div>
        {/each}
        <div class="input-group m-1">
          <input type="text" class="form-control" bind:value={newLink} />
          <IconButton color="success" icon="plus" tooltip="Add link" click={() => addLink()} />
        </div>
      </div>
    </div>
    <div class="row m-1">
      <div class="col-md">
        <h4>Pronouns</h4>
        {#each pronouns as _, index}
          <EditablePronouns
            bind:pronoun={pronouns[index]}
            moveUp={() => movePronoun(index, true)}
            moveDown={() => movePronoun(index, false)}
            remove={() => removePronoun(index)}
          />
        {/each}
        <div class="input-group m-1">
          <input
            type="text"
            class="form-control"
            placeholder="Full set (e.g. it/it/its/its/itself)"
            bind:value={newPronouns}
            required
          />
          <input
            type="text"
            class="form-control"
            placeholder="Optional display text (e.g. it/its)"
            bind:value={newPronounsDisplay}
          />
          <IconButton
            color="success"
            icon="plus"
            tooltip="Add pronouns"
            click={() => addPronouns()}
          />
        </div>
      </div>
    </div>
    <hr />
    <h4>
      Fields <Button on:click={() => (fields = [...fields, { name: "New field", entries: [] }])}>
        Add new field
      </Button>
    </h4>
  </div>
  <div class="grid gap-3">
    <div class="row row-cols-1 row-cols-md-2">
      {#each fields as _, index}
        <EditableField
          bind:field={fields[index]}
          deleteField={() => removeField(index)}
          moveField={(up) => moveField(index, up)}
        />
      {/each}
    </div>
  </div>
{/if}
