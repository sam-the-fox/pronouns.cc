<script lang="ts">
  import { flagURL, type APIError, type PrideFlag } from "$lib/api/entities";
  import { apiFetchClient } from "$lib/api/fetch";
  import { addToast } from "$lib/toast";
  import { Button, Input, Modal, ModalBody, ModalFooter, ModalHeader } from "@sveltestrap/sveltestrap";

  export let flag: PrideFlag;
  export let deleteFlag: (id: string) => Promise<void>;

  let error: APIError | null = null;

  let modalOpen = false;
  const toggleModal = () => (modalOpen = !modalOpen);

  let deleteModalOpen = false;
  const toggleDeleteModal = () => (deleteModalOpen = !deleteModalOpen);

  let name = flag.name;
  let description = flag.description;

  const updateFlag = async () => {
    try {
      const resp = await apiFetchClient<PrideFlag>(`/users/@me/flags/${flag.id_new}`, "PATCH", {
        name,
        description: description || null,
      });

      error = null;
      flag = resp;

      addToast({ header: "Updated flag", body: "Successfully updated flag!" });
      toggleModal();
    } catch (e) {
      error = e as APIError;
    }
  };
</script>

<Button outline class="m-1" on:click={toggleModal}>
  <img class="flag" src={flagURL(flag)} alt={flag.description ?? flag.name} />
  {flag.name}
</Button>

<Modal isOpen={modalOpen} toggle={toggleModal}>
  <ModalHeader toggle={toggleModal}>Edit {flag.name} flag</ModalHeader>
  <ModalBody>
    <p>
      <label for="name" class="form-label">Name</label>
      <Input id="name" bind:value={name} />
    </p>
    <p>
      <label for="description" class="form-label">Description</label>
      <textarea id="description" class="form-control" bind:value={description} />
    </p>
  </ModalBody>
  <ModalFooter>
    <Button color="danger" on:click={toggleDeleteModal}>Delete flag</Button>
    <Button disabled={!name} color="success" on:click={() => updateFlag()}>Edit flag</Button>
  </ModalFooter>
</Modal>

<Modal isOpen={deleteModalOpen} toggle={toggleDeleteModal}>
  <ModalHeader toggle={toggleDeleteModal}>Delete {flag.name} flag</ModalHeader>
  <ModalBody>
    Are you sure you want to delete the {flag.name} flag? <strong>This cannot be undone!</strong>
  </ModalBody>
  <ModalFooter>
    <Button color="danger" on:click={() => deleteFlag(flag.id_new)}>Delete flag</Button>
    <Button color="secondary" on:click={toggleDeleteModal}>Cancel</Button>
  </ModalFooter>
</Modal>

<style>
  .flag {
    height: 2rem;
    max-width: 200px;
    border-radius: 3px;
    margin-left: -5px;
  }

  textarea {
    height: 100px;
  }
</style>
