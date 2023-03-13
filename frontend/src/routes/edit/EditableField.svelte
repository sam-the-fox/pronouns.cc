<script lang="ts">
  import { WordStatus, type Field } from "$lib/api/entities";
  import IconButton from "$lib/components/IconButton.svelte";
  import { Button, Input, InputGroup, InputGroupText } from "sveltestrap";

  export let field: Field;
  export let deleteField: () => void;

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
    if (field.entries.length === 1) field.entries = [];
    else if (index === 0) field.entries = field.entries.slice(1);
    else if (index === field.entries.length - 1)
      field.entries = field.entries.slice(0, field.entries.length - 1);
    else
      field.entries = [...field.entries.slice(0, index - 1), ...field.entries.slice(0, index + 1)];
  };
</script>

<div>
  <InputGroup class="m-1">
    <InputGroupText>Name</InputGroupText>
    <Input bind:value={field.name} />
    <Button color="danger" on:click={() => deleteField()}>Delete field</Button>
  </InputGroup>
  {#each field.entries as _, index}
    <div class="input-group m-1">
      <IconButton
        color="secondary"
        icon="chevron-up"
        tooltip="Move entry up"
        click={() => moveEntry(index, true)}
      />
      <IconButton
        color="secondary"
        icon="chevron-down"
        tooltip="Move entry down"
        click={() => moveEntry(index, false)}
      />
      <input type="text" class="form-control" bind:value={field.entries[index].value} />
      <IconButton
        color="secondary"
        icon="heart-fill"
        tooltip="Favourite"
        click={() => (field.entries[index].status = WordStatus.Favourite)}
        active={field.entries[index].status === WordStatus.Favourite}
      />
      <IconButton
        color="secondary"
        icon="hand-thumbs-up"
        tooltip="Okay"
        click={() => (field.entries[index].status = WordStatus.Okay)}
        active={field.entries[index].status === WordStatus.Okay}
      />
      <IconButton
        color="secondary"
        icon="emoji-laughing"
        tooltip="Jokingly"
        click={() => (field.entries[index].status = WordStatus.Jokingly)}
        active={field.entries[index].status === WordStatus.Jokingly}
      />
      <IconButton
        color="secondary"
        icon="people"
        tooltip="Friends only"
        click={() => (field.entries[index].status = WordStatus.FriendsOnly)}
        active={field.entries[index].status === WordStatus.FriendsOnly}
      />
      <IconButton
        color="secondary"
        icon="hand-thumbs-down"
        tooltip="Avoid"
        click={() => (field.entries[index].status = WordStatus.Avoid)}
        active={field.entries[index].status === WordStatus.Avoid}
      />
      <IconButton
        color="danger"
        icon="trash3"
        tooltip="Remove entry"
        click={() => removeEntry(index)}
      />
    </div>
  {/each}

  <div class="input-group m-1">
    <input type="text" class="form-control" placeholder="New entry" bind:value={newEntry} />
    <IconButton color="success" icon="plus" tooltip="Add entry" click={() => addEntry()} />
  </div>
</div>
