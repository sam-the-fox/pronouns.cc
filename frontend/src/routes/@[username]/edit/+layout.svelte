<script lang="ts">
  import { setContext } from "svelte";
  import { writable } from "svelte/store";
  import type { LayoutData } from "./$types";
  import { Button, ButtonGroup, Icon, Nav, NavItem } from "sveltestrap";
  import type { MeUser, APIError } from "$lib/api/entities";
  import ErrorAlert from "$lib/components/ErrorAlert.svelte";
  import { addToast, delToast } from "$lib/toast";
  import { apiFetchClient } from "$lib/api/fetch";
  import { userStore } from "$lib/store";
  import ActiveLink from "$lib/components/ActiveLink.svelte";

  export let data: LayoutData;

  // The user context used in all other pages.
  // Avatar is explicitly set to null here as it holds the base64-encoded version of the new avatar later.
  const user = writable<MeUser>(structuredClone({ ...data.user, avatar: null }));
  const currentUser = writable<MeUser>(data.user);
  setContext("user", user);
  setContext("currentUser", currentUser);

  let error: APIError | null = null;

  const updateUser = async () => {
    const toastId = addToast({
      header: "Saving changes",
      body: "Saving changes, please wait...",
      duration: -1,
    });

    try {
      const resp = await apiFetchClient<MeUser>("/users/@me", "PATCH", {
        display_name: $user.display_name || "",
        avatar: $user.avatar,
        bio: $user.bio,
        links: $user.links,
        names: $user.names,
        pronouns: $user.pronouns,
        fields: $user.fields,
        member_title: $user.member_title || "",
        list_private: $user.list_private,
        timezone: $user.timezone || "",
        custom_preferences: $user.custom_preferences,
        flags: $user.flags.map((flag) => flag.id),
      });

      currentUser.set(resp);
      userStore.set(resp);
      localStorage.setItem("pronouns-user", JSON.stringify(resp));

      addToast({ header: "Success", body: "Successfully saved changes!" });

      user.update((_) => ({ ...resp, avatar: null }));
      error = null;
    } catch (e) {
      error = e as APIError;
    } finally {
      delToast(toastId);
    }
  };
</script>

<svelte:head>
  <title>Edit profile - pronouns.cc</title>
</svelte:head>

<h1>
  Edit profile
  <ButtonGroup>
    <Button color="secondary" href="/@{data.user.name}">
      <Icon name="chevron-left" />
      Back to your profile
    </Button>
    <Button color="success" on:click={() => updateUser()}>Save changes</Button>
  </ButtonGroup>
</h1>

{#if error}
  <ErrorAlert {error} />
{/if}

<Nav tabs>
  <NavItem><ActiveLink href="/@{$user.name}/edit">Names and avatar</ActiveLink></NavItem>
  <NavItem><ActiveLink href="/@{$user.name}/edit/bio">Bio</ActiveLink></NavItem>
  <NavItem><ActiveLink href="/@{$user.name}/edit/pronouns">Pronouns</ActiveLink></NavItem>
  <NavItem><ActiveLink href="/@{$user.name}/edit/fields">Fields</ActiveLink></NavItem>
  <NavItem><ActiveLink href="/@{$user.name}/edit/flags">Flags</ActiveLink></NavItem>
  <NavItem><ActiveLink href="/@{$user.name}/edit/links">Links</ActiveLink></NavItem>
  <NavItem><ActiveLink href="/@{$user.name}/edit/other">Preferences & other</ActiveLink></NavItem>
</Nav>

<div class="mt-3">
  <slot />
</div>
