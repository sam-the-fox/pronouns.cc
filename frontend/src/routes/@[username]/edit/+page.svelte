<script lang="ts">
  import { getContext } from "svelte";
  import type { Writable } from "svelte/store";
  import prettyBytes from "pretty-bytes";
  import { encode } from "base64-arraybuffer";
  import { FormGroup, Icon, Input } from "@sveltestrap/sveltestrap";
  import { userAvatars, type MeUser } from "$lib/api/entities";
  import FallbackImage from "$lib/components/FallbackImage.svelte";
  import EditableName from "$lib/components/edit/EditableName.svelte";
  import { addToast } from "$lib/toast";
  import IconButton from "$lib/components/IconButton.svelte";

  const MAX_AVATAR_BYTES = 1_000_000;

  const user = getContext<Writable<MeUser>>("user");
  const currentUser = getContext<Writable<MeUser>>("currentUser");

  // The list of avatar files uploaded by the user.
  // Only the first of these is ever used.
  let avatar_files: FileList | null;
  $: getAvatar(avatar_files).then((b64) => ($user.avatar = b64));

  // The variable for a new name being inputted by the user.
  let newName = "";

  const moveName = (index: number, up: boolean) => {
    if (up && index == 0) return;
    if (!up && index == $user.names.length - 1) return;

    const newIndex = up ? index - 1 : index + 1;

    const temp = $user.names[index];
    $user.names[index] = $user.names[newIndex];
    $user.names[newIndex] = temp;
    $user.names = [...$user.names];
  };

  const removeName = (index: number) => {
    $user.names.splice(index, 1);
    $user.names = [...$user.names];
  };

  const addName = (event: Event) => {
    event.preventDefault();

    $user.names = [...$user.names, { value: newName, status: "okay" }];
    newName = "";
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
</script>

<div class="row">
  <div class="col-md">
    <div class="row">
      <div class="col-md text-center">
        {#if $user.avatar === ""}
          <FallbackImage alt="Current avatar" urls={[]} width={200} />
        {:else if $user.avatar}
          <img
            width={200}
            height={200}
            src={$user.avatar}
            alt="New avatar"
            class="rounded-circle img-fluid"
          />
        {:else}
          <FallbackImage alt="Current avatar" urls={userAvatars($currentUser)} width={200} />
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
          <Icon name="info-circle-fill" aria-hidden /> Only PNG, JPEG, GIF, and WebP images can be used
          as avatars. Avatars cannot be larger than 1 MB, and animated avatars will be made static.
        </p>
        <p>
          <!-- svelte-ignore a11y-invalid-attribute -->
          <a href="" on:click={() => ($user.avatar = "")}>Remove avatar</a>
        </p>
      </div>
    </div>
  </div>
  <div class="col-md">
    <FormGroup floating label="Username">
      <Input bind:value={$user.name} readonly />
      <p class="text-muted mt-1">
        <Icon name="info-circle-fill" aria-hidden />
        You can change your username in
        <a href="/settings" class="text-reset">your settings</a>.
      </p>
    </FormGroup>
    <FormGroup floating label="Display name">
      <Input bind:value={$user.display_name} />
      <p class="text-muted mt-1">
        <Icon name="info-circle-fill" aria-hidden />
        Your display name is used in page titles and as a header.
      </p>
    </FormGroup>
  </div>
</div>
<div>
  <h4>Names</h4>
  {#each $user.names as _, index}
    <EditableName
      bind:value={$user.names[index].value}
      bind:status={$user.names[index].status}
      preferences={$user.custom_preferences}
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
