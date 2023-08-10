<script lang="ts">
  import { getContext } from "svelte";
  import type { Writable } from "svelte/store";
  import { Alert, Button, Icon } from "sveltestrap";

  import type { MeUser } from "$lib/api/entities";
  import EditableField from "$lib/components/edit/EditableField.svelte";

  const user = getContext<Writable<MeUser>>("user");

  const moveField = (index: number, up: boolean) => {
    if (up && index == 0) return;
    if (!up && index == $user.fields.length - 1) return;

    const newIndex = up ? index - 1 : index + 1;

    const temp = $user.fields[index];
    $user.fields[index] = $user.fields[newIndex];
    $user.fields[newIndex] = temp;
    $user.fields = [...$user.fields];
  };

  const removeField = (index: number) => {
    $user.fields.splice(index, 1);
    $user.fields = [...$user.fields];
  };
</script>

{#if $user.fields.length === 0}
  <Alert class="mt-3" color="secondary" fade={false}>
    Fields are extra categories you can add separate from names and pronouns.<br />
    For example, you could use them for gender terms, honorifics, or compliments.
  </Alert>
{/if}
<div class="grid gap-3">
  <div class="row row-cols-1 row-cols-md-2">
    {#each $user.fields as _, index}
      <EditableField
        bind:field={$user.fields[index]}
        preferences={$user.custom_preferences}
        deleteField={() => removeField(index)}
        moveField={(up) => moveField(index, up)}
      />
    {/each}
  </div>
</div>
<div>
  <Button on:click={() => ($user.fields = [...$user.fields, { name: "New field", entries: [] }])}>
    <Icon name="plus" aria-hidden /> Add new field
  </Button>
</div>
