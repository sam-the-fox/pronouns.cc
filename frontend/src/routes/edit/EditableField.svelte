<script lang="ts">
  import { WordStatus, type Field } from "$lib/api/entities";
  import IconButton from "$lib/components/IconButton.svelte";
  import { Button, Input, InputGroup, InputGroupText } from "sveltestrap";
  import FieldEntry from "./FieldEntry.svelte";

  export let field: Field;
  export let deleteField: () => void;
  export let moveField: (up: boolean) => void;

  let newEntry: string = "";

  const addEntry = () => {
    field.entries = [...field.entries, { value: newEntry, status: WordStatus.Okay }];
    newEntry = "";
  };

  const moveEntry = (index: number, up: boolean) => {
    if (up && index == 0) return;
    if (!up && index == field.entries.length - 1) return;

    const newIndex = up ? index - 1 : index + 1;

    const temp = field.entries[index];
    field.entries[index] = field.entries[newIndex];
    field.entries[newIndex] = temp;
  };

  const removeEntry = (index: number) => {
    field.entries.splice(index, 1);
    field.entries = [...field.entries];
  };
</script>

<div class="p-3">
  <h5>{field.name}</h5>
  <InputGroup class="m-1">
    <IconButton
      color="secondary"
      icon="chevron-left"
      tooltip="Move field left"
      click={() => moveField(true)}
    />
    <IconButton
      color="secondary"
      icon="chevron-right"
      tooltip="Move field right"
      click={() => moveField(false)}
    />
    <!-- <InputGroupText>Name</InputGroupText> -->
    <Input bind:value={field.name} />
    <Button color="danger" on:click={() => deleteField()}>Delete field</Button>
  </InputGroup>
  <hr />
  {#each field.entries as _, index}
    <FieldEntry
      bind:value={field.entries[index].value}
      bind:status={field.entries[index].status}
      moveUp={() => moveEntry(index, true)}
      moveDown={() => moveEntry(index, false)}
      remove={() => removeEntry(index)}
    />
  {/each}

  <div class="input-group m-1">
    <input type="text" class="form-control" placeholder="New entry" bind:value={newEntry} />
    <IconButton color="success" icon="plus" tooltip="Add entry" click={() => addEntry()} />
  </div>
</div>
