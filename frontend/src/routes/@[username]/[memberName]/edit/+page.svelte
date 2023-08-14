<script lang="ts">
  import { getContext } from "svelte";
  import type { Writable } from "svelte/store";
  import prettyBytes from "pretty-bytes";
  import { encode } from "base64-arraybuffer";
  import { FormGroup, Icon, Input } from "sveltestrap";
  import { memberAvatars, type Member } from "$lib/api/entities";
  import FallbackImage from "$lib/components/FallbackImage.svelte";
  import EditableName from "$lib/components/edit/EditableName.svelte";
  import { addToast } from "$lib/toast";
  import IconButton from "$lib/components/IconButton.svelte";
  import { memberNameRegex } from "$lib/api/regex";

  const MAX_AVATAR_BYTES = 1_000_000;

  const member = getContext<Writable<Member>>("member");
  const currentMember = getContext<Writable<Member>>("currentMember");

  // Whether the member's new name is valid.
  // This is also checked in +layout.svelte, because it's easier to do it twice.
  let memberNameValid = true;
  $: memberNameValid = memberNameRegex.test($member.name);

  // The list of avatar files uploaded.
  // Only the first of these is ever used.
  let avatar_files: FileList | null;
  $: getAvatar(avatar_files).then((b64) => ($member.avatar = b64));

  // The variable for a new name being inputted.
  let newName = "";

  const moveName = (index: number, up: boolean) => {
    if (up && index == 0) return;
    if (!up && index == $member.names.length - 1) return;

    const newIndex = up ? index - 1 : index + 1;

    const temp = $member.names[index];
    $member.names[index] = $member.names[newIndex];
    $member.names[newIndex] = temp;
    $member.names = [...$member.names];
  };

  const removeName = (index: number) => {
    $member.names.splice(index, 1);
    $member.names = [...$member.names];
  };

  const addName = (event: Event) => {
    event.preventDefault();

    $member.names = [...$member.names, { value: newName, status: "okay" }];
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
        {#if $member.avatar === ""}
          <FallbackImage alt="Current avatar" urls={[]} width={200} />
        {:else if $member.avatar}
          <img
            width={200}
            height={200}
            src={$member.avatar}
            alt="New avatar"
            class="rounded-circle img-fluid"
          />
        {:else}
          <FallbackImage alt="Current avatar" urls={memberAvatars($currentMember)} width={200} />
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
          <a href="" on:click={() => ($member.avatar = "")}>Remove avatar</a>
        </p>
      </div>
    </div>
  </div>
  <div class="col-md">
    <FormGroup floating label="Name">
      <Input bind:value={$member.name} />
      <p class="text-muted mt-1">
        <Icon name="info-circle-fill" aria-hidden />
        The member name is only used as part of the link to their profile page.
      </p>
    </FormGroup>
    {#if !memberNameValid}
      <p class="text-danger-emphasis mb-2">That member name is not valid.</p>
    {/if}
    <FormGroup floating label="Display name">
      <Input bind:value={$member.display_name} />
    </FormGroup>
    <p class="text-muted mt-1">
      <Icon name="info-circle-fill" aria-hidden />
      Your display name is used in page titles and as a header.
    </p>
  </div>
</div>
<div>
  <h4>Names</h4>
  {#each $member.names as _, index}
    <EditableName
      bind:value={$member.names[index].value}
      bind:status={$member.names[index].status}
      preferences={$member.user.custom_preferences}
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
