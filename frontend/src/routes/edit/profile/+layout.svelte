<script lang="ts">
  import { setContext } from "svelte";
  import { writable } from "svelte/store";
  import { page } from "$app/stores";
  import type { LayoutData } from "./$types";
  import { Button, ButtonGroup, Icon, Nav, NavItem, NavLink } from "sveltestrap";
  import type {
    FieldEntry,
    Field,
    MeUser,
    Pronoun,
    PrideFlag,
    CustomPreferences,
    APIError,
  } from "$lib/api/entities";
  import ErrorAlert from "$lib/components/ErrorAlert.svelte";
  import { addToast, delToast } from "$lib/toast";
  import { apiFetchClient } from "$lib/api/fetch";
  import { userStore } from "$lib/store";

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
  <NavItem>
    <NavLink href="/edit/profile" active={$page.url.pathname === "/edit/profile"}
      >Names and avatar</NavLink
    >
  </NavItem>
  <NavItem>
    <NavLink href="/edit/profile/bio" active={$page.url.pathname === "/edit/profile/bio"}
      >Bio</NavLink
    >
  </NavItem>
  <NavItem>
    <NavLink href="/edit/profile/pronouns" active={$page.url.pathname === "/edit/profile/pronouns"}
      >Pronouns</NavLink
    >
  </NavItem>
  <NavItem>
    <NavLink href="/edit/profile/fields" active={$page.url.pathname === "/edit/profile/fields"}
      >Fields</NavLink
    >
  </NavItem>
  <NavItem>
    <NavLink href="/edit/profile/flags" active={$page.url.pathname === "/edit/profile/flags"}
      >Flags</NavLink
    >
  </NavItem>
  <NavItem>
    <NavLink href="/edit/profile/links" active={$page.url.pathname === "/edit/profile/links"}
      >Links</NavLink
    >
  </NavItem>
  <NavItem>
    <NavLink href="/edit/profile/other" active={$page.url.pathname === "/edit/profile/other"}
      >Preferences & other</NavLink
    >
  </NavItem>
</Nav>

<div class="mt-3">
  <slot />
</div>
