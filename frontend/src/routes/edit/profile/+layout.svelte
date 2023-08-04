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
  let modified = false;
  $: modified = isModified($currentUser, $user);

  /** Returns whether or not the user in the store (referred to as newUser) is modified. */
  const isModified = (user: MeUser, newUser: MeUser) => {
    if ((newUser.bio || "") !== (user.bio || "")) return true;
    if ((newUser.display_name || "") !== (user.display_name || "")) return true;
    if ((newUser.member_title || "") !== (user.member_title || "")) return true;
    if (!linksEqual(newUser.links, user.links)) return true;
    if (!fieldsEqual(newUser.fields, user.fields)) return true;
    if (!flagsEqual(newUser.flags, user.flags)) return true;
    if (!namesEqual(newUser.names, user.names)) return true;
    if (!pronounsEqual(newUser.pronouns, user.pronouns)) return true;
    if (!customPreferencesEqual(newUser.custom_preferences, user.custom_preferences)) return true;
    if (newUser.avatar !== null) return true;
    if (newUser.list_private !== user.list_private) return true;
    if (newUser.timezone !== user.timezone) return true;

    return false;
  };

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

  // The individual functions called in isModified follow.
  const fieldsEqual = (arr1: Field[], arr2: Field[]) => {
    if (arr1?.length !== arr2?.length) return false;
    if (!arr1.every((_, i) => arr1[i].entries.length === arr2[i].entries.length)) return false;
    if (!arr1.every((_, i) => arr1[i].name === arr2[i].name)) return false;

    return arr1.every((_, i) =>
      arr1[i].entries.every(
        (entry, j) =>
          entry.value === arr2[i].entries[j].value && entry.status === arr2[i].entries[j].status,
      ),
    );
  };

  const namesEqual = (arr1: FieldEntry[], arr2: FieldEntry[]) => {
    console.log("new:", JSON.stringify(arr1));
    console.log("current:", JSON.stringify(arr2));

    if (arr1?.length !== arr2?.length) return false;
    if (!arr1.every((_, i) => arr1[i].value === arr2[i].value)) return false;
    if (!arr1.every((_, i) => arr1[i].status === arr2[i].status)) return false;

    return true;
  };

  const pronounsEqual = (arr1: Pronoun[], arr2: Pronoun[]) => {
    if (arr1?.length !== arr2?.length) return false;
    if (!arr1.every((_, i) => arr1[i].pronouns === arr2[i].pronouns)) return false;
    if (!arr1.every((_, i) => arr1[i].display_text === arr2[i].display_text)) return false;
    if (!arr1.every((_, i) => arr1[i].status === arr2[i].status)) return false;

    return true;
  };

  const linksEqual = (arr1: string[], arr2: string[]) => {
    if (arr1.length !== arr2.length) return false;
    return arr1.every((_, i) => arr1[i] === arr2[i]);
  };

  const flagsEqual = (arr1: PrideFlag[], arr2: PrideFlag[]) => {
    if (arr1.length !== arr2.length) return false;
    return arr1.every((_, i) => arr1[i].id === arr2[i].id);
  };

  const customPreferencesEqual = (obj1: CustomPreferences, obj2: CustomPreferences) => {
    if (Object.keys(obj2).some((key) => !(key in obj1))) return false;

    return Object.keys(obj1)
      .map((key) => {
        if (!(key in obj2)) return false;
        return (
          obj1[key].icon === obj2[key].icon &&
          obj1[key].tooltip === obj2[key].tooltip &&
          obj1[key].favourite === obj2[key].favourite &&
          obj1[key].muted === obj2[key].muted &&
          obj1[key].size === obj2[key].size
        );
      })
      .every((entry) => entry);
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
    {#if modified}
      <Button color="success" on:click={() => updateUser()}>Save changes</Button>
    {/if}
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
