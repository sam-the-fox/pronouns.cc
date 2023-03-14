<script lang="ts">
  import { goto } from "$app/navigation";
  import {
    MAX_DESCRIPTION_LENGTH,
    userAvatars,
    WordStatus,
    type APIError,
    type Field,
    type FieldEntry,
    type MeUser,
    type Pronoun,
  } from "$lib/api/entities";
  import FallbackImage from "$lib/components/FallbackImage.svelte";
  import { userStore } from "$lib/store";
  import {
    Alert,
    Button,
    ButtonDropdown,
    DropdownToggle,
    FormGroup,
    Icon,
    Input,
  } from "sveltestrap";
  import { encode } from "base64-arraybuffer";
  import { apiFetchClient } from "$lib/api/fetch";
  import IconButton from "$lib/components/IconButton.svelte";
  import EditableField from "../EditableField.svelte";
  import EditableName from "../EditableName.svelte";
  import EditablePronouns from "../EditablePronouns.svelte";
  import ErrorAlert from "$lib/components/ErrorAlert.svelte";
  import { addToast, delToast } from "$lib/toast";

  const MAX_AVATAR_BYTES = 1_000_000;

  if (!$userStore) {
    goto("/");
  }

  let error: APIError | null = null;

  let bio: string = $userStore?.bio || "";
  let display_name: string = $userStore?.display_name || "";
  let links: string[] = $userStore ? window.structuredClone($userStore.links) : [];
  let names: FieldEntry[] = $userStore ? window.structuredClone($userStore.names) : [];
  let pronouns: Pronoun[] = $userStore ? window.structuredClone($userStore.pronouns) : [];
  let fields: Field[] = $userStore ? window.structuredClone($userStore.fields) : [];

  let avatar: string | null;
  let avatar_files: FileList | null;

  let newName = "";
  let newPronouns = "";
  let newPronounsDisplay = "";
  let newLink = "";

  let modified = false;

  $: redirectIfNoAuth($userStore);
  $: modified = isModified(bio, display_name, links, names, pronouns, fields, avatar);
  $: getAvatar(avatar_files).then((b64) => (avatar = b64));

  const redirectIfNoAuth = (user: MeUser | null) => {
    if (!user) {
      goto("/");
    }
  };

  const isModified = (
    bio: string,
    display_name: string,
    links: string[],
    names: FieldEntry[],
    pronouns: Pronoun[],
    fields: Field[],
    avatar: string | null,
  ) => {
    if (!$userStore) return false;

    if (bio !== $userStore.bio) return true;
    if (display_name !== $userStore.display_name) return true;
    if (!linksEqual(links, $userStore.links)) return true;
    if (!fieldsEqual(fields, $userStore.fields)) return true;
    if (!namesEqual(names, $userStore.names)) return true;
    if (!pronounsEqual(pronouns, $userStore.pronouns)) return true;
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

    const uri = `data:${dataTypeFromFilename(list[0].name)};base64,${base64}`;

    return uri;
  };

  const dataTypeFromFilename = (filename: string) => {
    if (filename.endsWith(".webp")) {
      return "image/webp";
    } else if (filename.endsWith(".jpg") || filename.endsWith(".jpeg")) {
      return "image/jpeg";
    } else if (filename.endsWith(".png")) {
      return "image/png";
    } else if (filename.endsWith(".gif")) {
      return "image/gif";
    } else {
      return "application/octet-stream";
    }
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
      });

      userStore.set(resp);
      localStorage.setItem("pronouns-user", JSON.stringify(resp));

      addToast({ header: "Success", body: "Successfully saved changes!" });

      avatar = null;
      error = null;
      modified = false;
    } catch (e) {
      error = e as APIError;
    } finally {
      delToast(toastId);
    }
  };
</script>

<svelte:head>
  <title>Edit profile - pronouns.cc</title>
</svelte:head>

<h1>
  Edit profile
  {#if modified}
    <Button color="success" on:click={() => updateUser()}>Save changes</Button>
  {/if}
</h1>

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
          <div class="col-md">
            {#if avatar}
              <img width={200} src={avatar} alt="New avatar" class="rounded-circle img-fluid" />
            {:else}
              <FallbackImage alt="Current avatar" urls={userAvatars($userStore)} width={200} />
            {/if}
          </div>
          <div class="col-md">
            <input class="form-control" id="avatar" type="file" bind:files={avatar_files} />
          </div>
        </div>
      </div>
      <div class="col-md">
        <FormGroup floating label="Display name">
          <Input bind:value={display_name} />
        </FormGroup>
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
        <EditableField bind:field={fields[index]} deleteField={() => removeField(index)} />
      {/each}
    </div>
  </div>
{/if}
