<script lang="ts">
  import { getContext } from "svelte";
  import type { Writable } from "svelte/store";
  import { Alert, Button, Icon } from "sveltestrap";

  import type { Member } from "$lib/api/entities";
  import EditableField from "$lib/components/edit/EditableField.svelte";

  const member = getContext<Writable<Member>>("member");

  const moveField = (index: number, up: boolean) => {
    if (up && index == 0) return;
    if (!up && index == $member.fields.length - 1) return;

    const newIndex = up ? index - 1 : index + 1;

    const temp = $member.fields[index];
    $member.fields[index] = $member.fields[newIndex];
    $member.fields[newIndex] = temp;
    $member.fields = [...$member.fields];
  };

  const removeField = (index: number) => {
    $member.fields.splice(index, 1);
    $member.fields = [...$member.fields];
  };
</script>

{#if $member.fields.length === 0}
  <Alert class="mt-3" color="secondary" fade={false}>
    Fields are extra categories you can add separate from names and pronouns.<br />
    For example, you could use them for gender terms, honorifics, or compliments.
  </Alert>
{/if}
<div class="grid gap-3">
  <div class="row row-cols-1 row-cols-md-2">
    {#each $member.fields as _, index}
      <EditableField
        bind:field={$member.fields[index]}
        preferences={$member.user.custom_preferences}
        deleteField={() => removeField(index)}
        moveField={(up) => moveField(index, up)}
      />
    {/each}
  </div>
</div>
<div>
  <Button
    on:click={() => ($member.fields = [...$member.fields, { name: "New field", entries: [] }])}
  >
    <Icon name="plus" aria-hidden /> Add new field
  </Button>
</div>
