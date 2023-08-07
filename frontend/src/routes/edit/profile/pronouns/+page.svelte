<script lang="ts">
  import { getContext } from "svelte";
  import type { Writable } from "svelte/store";
  import type { MeUser } from "$lib/api/entities";
  import { Button, Icon, Popover } from "sveltestrap";
  import EditablePronouns from "../../EditablePronouns.svelte";
  import IconButton from "$lib/components/IconButton.svelte";
  import type { PageData } from "./$types";

  export let data: PageData;

  const user = getContext<Writable<MeUser>>("user");
  let newPronouns = "";

  const movePronoun = (index: number, up: boolean) => {
    if (up && index == 0) return;
    if (!up && index == $user.pronouns.length - 1) return;

    const newIndex = up ? index - 1 : index + 1;

    const temp = $user.pronouns[index];
    $user.pronouns[index] = $user.pronouns[newIndex];
    $user.pronouns[newIndex] = temp;
    $user.pronouns = [...$user.pronouns];
  };

  const addPronouns = (event: Event) => {
    event.preventDefault();

    if (newPronouns in data.pronouns) {
      const fullSet = data.pronouns[newPronouns];
      $user.pronouns = [
        ...$user.pronouns,
        {
          pronouns: fullSet.pronouns.join("/"),
          display_text: fullSet.display || null,
          status: "okay",
        },
      ];
    } else {
      $user.pronouns = [
        ...$user.pronouns,
        { pronouns: newPronouns, display_text: null, status: "okay" },
      ];
    }
    newPronouns = "";
  };

  const removePronoun = (index: number) => {
    $user.pronouns.splice(index, 1);
    $user.pronouns = [...$user.pronouns];
  };
</script>

{#each $user.pronouns as _, index}
  <EditablePronouns
    bind:pronoun={$user.pronouns[index]}
    preferences={$user.custom_preferences}
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
    For common pronouns, the short form (e.g. "she/her" or "he/him") is enough; for less common
    pronouns, you will have to use all five forms (e.g. "ce/cir/cir/cirs/cirself").
  </Popover>
</form>
