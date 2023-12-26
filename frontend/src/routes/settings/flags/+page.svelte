<script lang="ts">
  import { MAX_FLAGS, type APIError, type PrideFlag } from "$lib/api/entities";
  import { Button, Icon, Input, Modal, ModalBody, ModalFooter, ModalHeader } from "@sveltestrap/sveltestrap";
  import type { PageData } from "./$types";
  import Flag from "./Flag.svelte";
  import prettyBytes from "pretty-bytes";
  import { addToast } from "$lib/toast";
  import { encode } from "base64-arraybuffer";
  import unknownFlag from "./unknown_flag.png";
  import { apiFetchClient, fastFetchClient } from "$lib/api/fetch";
  import ErrorAlert from "$lib/components/ErrorAlert.svelte";

  const MAX_FLAG_BYTES = 500_000;

  export let data: PageData;

  let search = "";
  let error: APIError | null = null;

  let filtered: PrideFlag[];
  $: filtered = filterFlags(search, data.flags);

  const filterFlags = (search: string, flags: PrideFlag[]) => {
    return search
      ? flags.filter((flag) => flag.name.toLocaleLowerCase().includes(search.toLocaleLowerCase()))
      : flags;
  };

  // NEW FLAG UPLOADING CODE
  let modalOpen = false;
  const toggleModal = () => (modalOpen = !modalOpen);
  let canUpload: boolean;
  $: canUpload = !!(newFlag && newName);

  let newFlag: string | null;
  let flagFiles: FileList | null;
  $: getFlag(flagFiles).then((b64) => (newFlag = b64));

  let newName = "";
  let newDescription = "";

  const getFlag = async (list: FileList | null) => {
    if (!list || list.length === 0) return null;
    if (list[0].size > MAX_FLAG_BYTES) {
      addToast({
        header: "Flag too large",
        body: `This flag file is too large, please resize it (maximum is ${prettyBytes(
          MAX_FLAG_BYTES,
        )}, the file you tried to upload is ${prettyBytes(list[0].size)})`,
      });
      return null;
    }

    const buffer = await list[0].arrayBuffer();
    const base64 = encode(buffer);

    const uri = `data:${list[0].type};base64,${base64}`;

    return uri;
  };

  const uploadFlag = async () => {
    try {
      const resp = await apiFetchClient<PrideFlag>("/users/@me/flags", "POST", {
        flag: newFlag,
        name: newName,
        description: newDescription || null,
      });

      error = null;
      data.flags.push(resp);
      data.flags.sort((a, b) => a.name.localeCompare(b.name));
      data.flags = [...data.flags];

      // reset flag
      newFlag = null;
      newName = "";
      newDescription = "";

      addToast({ header: "Uploaded flag", body: "Successfully uploaded flag!" });
      toggleModal();
    } catch (e) {
      error = e as APIError;
    }
  };

  // DELETE FLAG CODE
  const deleteFlag = async (id: string) => {
    try {
      await fastFetchClient(`/users/@me/flags/${id}`, "DELETE");

      error = null;

      addToast({ header: "Deleted flag", body: "Successfully deleted flag!" });
      data.flags = data.flags.filter((entry) => entry.id_new !== id);
    } catch (e) {
      error = e as APIError;
    }
  };
</script>

<h1>Pride flags ({data.flags.length}/{MAX_FLAGS})</h1>

<p>
  You can upload pride flags to use on your profiles here. Flags you upload here will <em>not</em> automatically
  show up on your profile.
</p>

<div class="input-group">
  <Input placeholder="Filter flags" bind:value={search} disabled={data.flags.length === 0} />
  <Button color="success" on:click={toggleModal}>
    <Icon name="upload" aria-hidden /> Upload flag
  </Button>
</div>

<div class="p-2">
  {#each filtered as flag (flag.id)}
    <Flag bind:flag {deleteFlag} />
  {:else}
    {#if data.flags.length === 0}
      You haven't uploaded any flags yet, press the button above to do so.
    {:else}
      There are no flags matching your search <strong>{search}</strong>.
    {/if}
  {/each}
</div>

<Modal isOpen={modalOpen} toggle={toggleModal}>
  <ModalHeader toggle={toggleModal}>Upload flag</ModalHeader>
  <ModalBody>
    {#if error}
      <ErrorAlert {error} />
    {/if}
    <div class="d-flex align-items-center">
      <img src={newFlag || unknownFlag} alt="New flag" class="flag m-1" />
      <input
        class="form-control"
        id="flag-file"
        type="file"
        bind:files={flagFiles}
        accept="image/png, image/jpeg, image/gif, image/webp, image/svg+xml"
      />
    </div>
    <p class="text-muted mt-2">
      <Icon name="info-circle-fill" aria-hidden /> Only PNG, JPEG, GIF, and WebP images can be uploaded
      as flags. The file cannot be larger than 512 kilobytes.
    </p>
    <p>
      <label for="newName" class="form-label">Name</label>
      <Input id="newName" bind:value={newName} />
    </p>
    <p class="text-muted">
      <Icon name="info-circle-fill" aria-hidden /> This name will be shown beside the flag.
    </p>
    <p>
      <label for="description" class="form-label">Description</label>
      <textarea id="description" class="form-control" bind:value={newDescription} />
    </p>
    <p class="text-muted">
      <Icon name="info-circle-fill" aria-hidden /> This text will be used as the alt text of the flag
      image, and will also be shown on hover. Optional, but <strong>strongly recommended</strong> as
      it improves accessibility.
    </p>
  </ModalBody>
  <ModalFooter>
    <Button disabled={!canUpload} color="success" on:click={() => uploadFlag()}>Upload flag</Button>
  </ModalFooter>
</Modal>

<style>
  .flag {
    height: 2rem;
    max-width: 200px;
    border-radius: 3px;
  }

  textarea {
    height: 100px;
  }
</style>
