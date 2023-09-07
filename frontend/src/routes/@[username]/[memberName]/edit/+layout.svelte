<script lang="ts">
  import type { APIError, Member, Pronoun } from "$lib/api/entities";
  import { setContext } from "svelte";
  import { writable } from "svelte/store";
  import type { LayoutData } from "./$types";
  import { addToast, delToast } from "$lib/toast";
  import { apiFetchClient, fastFetchClient } from "$lib/api/fetch";
  import { Button, ButtonGroup, Modal, ModalBody, ModalFooter, Nav, NavItem } from "sveltestrap";
  import { goto } from "$app/navigation";
  import ErrorAlert from "$lib/components/ErrorAlert.svelte";
  import IconButton from "$lib/components/IconButton.svelte";
  import ActiveLink from "$lib/components/ActiveLink.svelte";
  import { memberNameRegex } from "$lib/api/regex";

  export let data: LayoutData;

  const member = writable<Member>(structuredClone({ ...data.member, avatar: null }));
  const currentMember = writable(data.member);
  setContext("member", member);
  setContext("currentMember", currentMember);

  // Whether the member's new name is valid.
  // This is also checked in +page.svelte, because it's easier to do it twice.
  let memberNameValid = true;
  $: memberNameValid = memberNameRegex.test($member.name);

  let error: APIError | null = null;

  // Delete member code
  let deleteOpen = false;
  const toggleDeleteOpen = () => (deleteOpen = !deleteOpen);
  let deleteName = "";
  let deleteError: APIError | null = null;

  const deleteMember = async () => {
    try {
      await fastFetchClient(`/members/${data.member.id_new}`, "DELETE");

      toggleDeleteOpen();
      addToast({
        header: "Deleted member",
        body: `Successfully deleted member ${data.member.name}!`,
      });
      goto(`/@${data.member.user.name}`);
    } catch (e) {
      deleteName = "";
      deleteError = e as APIError;
    }
  };

  let deleteModalPronoun = "the member's";
  $: deleteModalPronoun = updateModalPronoun($member.pronouns);

  const updateModalPronoun = (pronouns: Pronoun[]) => {
    const filtered = pronouns.filter((entry) => entry.status === "favourite");
    if (filtered.length < 1) return "the member's";

    const split = filtered[0].pronouns.split("/");
    if (split.length !== 5) return "the member's";
    return split[2];
  };

  const updateMember = async () => {
    const toastId = addToast({
      header: "Saving changes",
      body: "Saving changes, please wait...",
      duration: -1,
    });

    try {
      const resp = await apiFetchClient<Member>(`/members/${data.member.id_new}`, "PATCH", {
        name: $member.name,
        display_name: $member.display_name,
        avatar: $member.avatar,
        bio: $member.bio,
        links: $member.links,
        names: $member.names,
        pronouns: $member.pronouns,
        fields: $member.fields,
        flags: $member.flags.map((flag) => flag.id),
        unlisted: $member.unlisted,
      });

      addToast({ header: "Success", body: "Successfully saved changes!" });

      data.member = resp;
      $member.avatar = null;
      error = null;
    } catch (e) {
      error = e as APIError;
    } finally {
      delToast(toastId);
    }
  };
</script>

<svelte:head>
  <title>Edit member profile - pronouns.cc</title>
</svelte:head>

<h1>
  Edit profile
  <ButtonGroup>
    <IconButton
      color="secondary"
      icon="chevron-left"
      href="/@{$member.user.name}/{data.member.name}"
      tooltip="Back to member"
    />
    <Button color="success" on:click={() => updateMember()} disabled={!memberNameValid}>
      Save changes
    </Button>
    <Button color="danger" on:click={toggleDeleteOpen}
      >Delete {data.member.display_name ?? data.member.name}</Button
    >
  </ButtonGroup>
</h1>

{#if error}
  <ErrorAlert {error} />
{/if}

<Nav tabs>
  <NavItem>
    <ActiveLink href="/@{$member.user.name}/{$member.name}/edit">Names and avatar</ActiveLink>
  </NavItem>
  <NavItem>
    <ActiveLink href="/@{$member.user.name}/{$member.name}/edit/bio">Bio</ActiveLink>
  </NavItem>
  <NavItem>
    <ActiveLink href="/@{$member.user.name}/{$member.name}/edit/pronouns">Pronouns</ActiveLink>
  </NavItem>
  <NavItem>
    <ActiveLink href="/@{$member.user.name}/{$member.name}/edit/fields">Fields</ActiveLink>
  </NavItem>
  <NavItem>
    <ActiveLink href="/@{$member.user.name}/{$member.name}/edit/flags">Flags</ActiveLink>
  </NavItem>
  <NavItem>
    <ActiveLink href="/@{$member.user.name}/{$member.name}/edit/links">Links</ActiveLink>
  </NavItem>
  <NavItem>
    <ActiveLink href="/@{$member.user.name}/{$member.name}/edit/other">Other</ActiveLink>
  </NavItem>
</Nav>

<div class="mt-3">
  <slot />
</div>

<Modal header="Delete member" isOpen={deleteOpen} toggle={toggleDeleteOpen}>
  <ModalBody>
    <p>
      If you want to delete this member, type {deleteModalPronoun} name (<code>{$member.name}</code
      >) below:
    </p>
    <p>
      <input type="text" class="form-control" bind:value={deleteName} />
    </p>
    {#if deleteError}
      <ErrorAlert error={deleteError} />
    {/if}
  </ModalBody>
  <ModalFooter>
    <Button color="danger" disabled={deleteName !== $member.name} on:click={deleteMember}>
      Delete member
    </Button>
    <Button color="secondary" on:click={toggleDeleteOpen}>Cancel</Button>
  </ModalFooter>
</Modal>
