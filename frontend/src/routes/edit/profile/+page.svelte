<script lang="ts">
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
  import { Button, ButtonGroup, FormGroup, Icon, Input, Popover } from "sveltestrap";
  import { encode } from "base64-arraybuffer";
  import { apiFetchClient } from "$lib/api/fetch";
  import IconButton from "$lib/components/IconButton.svelte";
  import EditableField from "../EditableField.svelte";
  import EditableName from "../EditableName.svelte";
  import EditablePronouns from "../EditablePronouns.svelte";
  import ErrorAlert from "$lib/components/ErrorAlert.svelte";
  import { addToast, delToast } from "$lib/toast";
  import type { PageData } from "./$types";

  const MAX_AVATAR_BYTES = 1_000_000;

  export let data: PageData;

  let error: APIError | null = null;

  let bio: string = data.user.bio || "";
  let display_name: string = data.user.display_name || "";
  let links: string[] = window.structuredClone(data.user.links);
  let names: FieldEntry[] = window.structuredClone(data.user.names);
  let pronouns: Pronoun[] = window.structuredClone(data.user.pronouns);
  let fields: Field[] = window.structuredClone(data.user.fields);

  let avatar: string | null;
  let avatar_files: FileList | null;

  let newName = "";
  let newPronouns = "";
  let newLink = "";

  let modified = false;

  $: modified = isModified(bio, display_name, links, names, pronouns, fields, avatar);
  $: getAvatar(avatar_files).then((b64) => (avatar = b64));

  const isModified = (
    bio: string,
    display_name: string,
    links: string[],
    names: FieldEntry[],
    pronouns: Pronoun[],
    fields: Field[],
    avatar: string | null,
  ) => {
    if (bio !== data.user.bio) return true;
    if (display_name !== data.user.display_name) return true;
    if (!linksEqual(links, data.user.links)) return true;
    if (!fieldsEqual(fields, data.user.fields)) return true;
    if (!namesEqual(names, data.user.names)) return true;
    if (!pronounsEqual(pronouns, data.user.pronouns)) return true;
    if (avatar !== null) return true;

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

  const getAvatar = async (list: FileList | null) => {
    if (!list || list.length === 0) return null;
    if (list[0].size > MAX_AVATAR_BYTES) return null;

    const buffer = await list[0].arrayBuffer();
    const base64 = encode(buffer);

    const uri = `data:${list[0].type};base64,${base64}`;
    console.log(uri.slice(0, 128));

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

  const addName = (event: Event) => {
    event.preventDefault();

    names = [...names, { value: newName, status: WordStatus.Okay }];
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
          status: WordStatus.Okay,
        },
      ];
    } else {
      pronouns = [
        ...pronouns,
        { pronouns: newPronouns, display_text: null, status: WordStatus.Okay },
      ];
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

      data.user = resp;
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

<div class="grid">
  <div class="row m-1">
    <div class="col-md">
      <h4>Avatar</h4>
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
        <div class="col-md mt-2">
          <input
            class="form-control"
            id="avatar"
            type="file"
            bind:files={avatar_files}
            accept="image/png, image/jpeg, image/gif, image/webp"
          />
          <p class="text-muted mt-3">
            <Icon name="info-circle-fill" aria-hidden /> Only PNG, JPEG, GIF, and WebP images can be
            used as avatars. Avatars cannot be larger than 1 MB, and animated avatars will be made static.
          </p>
          <p>
            <a href="" on:click={() => (avatar = "")}>Remove avatar</a>
          </p>
        </div>
      </div>
    </div>
    <div class="col-md">
      <FormGroup floating label="Display name">
        <Input bind:value={display_name} />
      </FormGroup>
      <div>
        <div class="form">
          <label for="bio"><strong>Bio ({bio.length}/{MAX_DESCRIPTION_LENGTH})</strong></label>
          <textarea class="form-control" style="height: 200px;" id="bio" bind:value={bio} />
        </div>
        <p class="text-muted mt-3">
          <Icon name="info-circle-fill" aria-hidden /> Your bio supports limited
          <a
            class="text-reset"
            href="https://commonmark.org/help/"
            target="_blank"
            rel="noopener noreferrer">Markdown</a
          >.
        </p>
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
      <form class="input-group m-1" on:submit={addName}>
        <input type="text" class="form-control" bind:value={newName} />
        <IconButton type="submit" color="success" icon="plus" tooltip="Add name" />
      </form>
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
      <form class="input-group m-1" on:submit={addLink}>
        <input type="text" class="form-control" bind:value={newLink} />
        <IconButton type="submit" color="success" icon="plus" tooltip="Add link" />
      </form>
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
