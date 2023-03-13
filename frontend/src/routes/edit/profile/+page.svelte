<script lang="ts">
  import { goto } from "$app/navigation";
  import {
    WordStatus,
    type APIError,
    type Field,
    type FieldEntry,
    type MeUser,
    type Pronoun,
  } from "$lib/api/entities";
  import FallbackImage from "$lib/components/FallbackImage.svelte";
  import { userStore } from "$lib/store";
  import { Alert, Button, FormGroup, Icon, Input } from "sveltestrap";
  import { encode } from "base64-arraybuffer";
  import { apiFetchClient } from "$lib/api/fetch";

  const MAX_AVATAR_BYTES = 1_000_000;

  if (!$userStore) {
    goto("/");
  }

  let error: APIError | null = null;

  let bio: string = $userStore?.bio || "";
  let display_name: string = $userStore?.display_name || "";
  let names: FieldEntry[] = $userStore ? window.structuredClone($userStore.names) : [];
  let pronouns: Pronoun[] = $userStore ? window.structuredClone($userStore.pronouns) : [];
  let fields: Field[] = $userStore ? window.structuredClone($userStore.fields) : [];

  let avatar: string | null;
  let avatar_files: FileList | null;

  let modified = false;

  $: redirectIfNoAuth($userStore);
  $: modified = isModified(bio, display_name, names, pronouns, fields);
  $: getAvatar(avatar_files).then((b64) => (avatar = b64));

  const redirectIfNoAuth = (user: MeUser | null) => {
    if (!user) {
      goto("/");
    }
  };

  const isModified = (
    bio: string,
    display_name: string,
    names: FieldEntry[],
    pronouns: Pronoun[],
    fields: Field[],
  ) => {
    if (!$userStore) return false;

    if (bio !== $userStore.bio) return true;
    if (display_name !== $userStore.display_name) return true;
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

  const updateUser = async () => {
    try {
      const resp = await apiFetchClient<MeUser>("/users/@me", "PATCH", {
        display_name,
        avatar,
        bio,
        names,
        pronouns,
        fields,
      });

      userStore.set(resp);
      localStorage.setItem("pronouns-user", JSON.stringify(resp));

      avatar = null;
      error = null;
      modified = false;
    } catch (e) {
      error = e as APIError;
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
  <Alert color="danger" fade={false}>
    <h4 class="alert-header">An error occurred</h4>
    <p>
      <b>{error.code}</b>: {error.message}
      {#if error.details}{error.details}{/if}
    </p>
  </Alert>
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
              <FallbackImage alt="Current avatar" urls={$userStore.avatar_urls} width={200} />
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
          <FormGroup floating label="Bio ({bio.length}/1000)">
            <textarea style="min-height: 100px;" class="form-control" bind:value={bio} />
          </FormGroup>
        </div>
      </div>
    </div>
    <div class="row m-1">
      <div class="col-md">
        <h4>Names</h4>
        {#each names as _, index}
          <div class="input-group m-1">
            <Button color="secondary" on:click={() => moveName(index, true)}>
              <Icon name="chevron-up" />
            </Button>
            <Button color="secondary" on:click={() => moveName(index, false)}>
              <Icon name="chevron-down" />
            </Button>
            <input type="text" class="form-control" bind:value={names[index].value} />
            <Button
              color="secondary"
              on:click={() => (names[index].status = WordStatus.Favourite)}
              active={names[index].status === WordStatus.Favourite}
            >
              <Icon name="heart-fill" />
            </Button>
            <Button
              color="secondary"
              on:click={() => (names[index].status = WordStatus.Okay)}
              active={names[index].status === WordStatus.Okay}
            >
              <Icon name="hand-thumbs-up" />
            </Button>
            <Button
              color="secondary"
              on:click={() => (names[index].status = WordStatus.Jokingly)}
              active={names[index].status === WordStatus.Jokingly}
            >
              <Icon name="emoji-laughing" />
            </Button>
            <Button
              color="secondary"
              on:click={() => (names[index].status = WordStatus.FriendsOnly)}
              active={names[index].status === WordStatus.FriendsOnly}
            >
              <Icon name="people" />
            </Button>
            <Button
              color="secondary"
              on:click={() => (names[index].status = WordStatus.Avoid)}
              active={names[index].status === WordStatus.Avoid}
            >
              <Icon name="hand-thumbs-down" />
            </Button>
            <Button color="danger">
              <Icon name="trash3" />
            </Button>
          </div>
        {/each}
      </div>
    </div>
    <div class="row m-1">
      <div class="col-md">
        <h4>Pronouns</h4>
        {#each pronouns as _, index}
          <div class="input-group m-1">
            <Button color="secondary" on:click={() => movePronoun(index, true)}>
              <Icon name="chevron-up" />
            </Button>
            <Button color="secondary" on:click={() => movePronoun(index, false)}>
              <Icon name="chevron-down" />
            </Button>
            <input type="text" class="form-control" bind:value={pronouns[index].pronouns} />
            <input type="text" class="form-control" bind:value={pronouns[index].display_text} />
            <Button
              color="secondary"
              on:click={() => (pronouns[index].status = WordStatus.Favourite)}
              active={pronouns[index].status === WordStatus.Favourite}
            >
              <Icon name="heart-fill" />
            </Button>
            <Button
              color="secondary"
              on:click={() => (pronouns[index].status = WordStatus.Okay)}
              active={pronouns[index].status === WordStatus.Okay}
            >
              <Icon name="hand-thumbs-up" />
            </Button>
            <Button
              color="secondary"
              on:click={() => (pronouns[index].status = WordStatus.Jokingly)}
              active={pronouns[index].status === WordStatus.Jokingly}
            >
              <Icon name="emoji-laughing" />
            </Button>
            <Button
              color="secondary"
              on:click={() => (pronouns[index].status = WordStatus.FriendsOnly)}
              active={pronouns[index].status === WordStatus.FriendsOnly}
            >
              <Icon name="people" />
            </Button>
            <Button
              color="secondary"
              on:click={() => (pronouns[index].status = WordStatus.Avoid)}
              active={pronouns[index].status === WordStatus.Avoid}
            >
              <Icon name="hand-thumbs-down" />
            </Button>
            <Button color="danger">
              <Icon name="trash3" />
            </Button>
          </div>
        {/each}
      </div>
    </div>
  </div>
{/if}
