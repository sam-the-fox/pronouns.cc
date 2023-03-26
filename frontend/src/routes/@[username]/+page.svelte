<script lang="ts">
  import type { PageData } from "./$types";

  import {
    Alert,
    Button,
    ButtonGroup,
    Icon,
    Input,
    Modal,
    ModalBody,
    ModalFooter,
  } from "sveltestrap";
  import FieldCard from "$lib/components/FieldCard.svelte";
  import StatusIcon from "$lib/components/StatusIcon.svelte";
  import PronounLink from "$lib/components/PronounLink.svelte";
  import PartialMemberCard from "$lib/components/PartialMemberCard.svelte";
  import FallbackImage from "$lib/components/FallbackImage.svelte";
  import { userStore } from "$lib/store";
  import {
    MAX_MEMBERS,
    pronounDisplay,
    userAvatars,
    WordStatus,
    type APIError,
    type Member,
    type PartialMember,
  } from "$lib/api/entities";
  import { PUBLIC_BASE_URL } from "$env/static/public";
  import { apiFetchClient } from "$lib/api/fetch";
  import ErrorAlert from "$lib/components/ErrorAlert.svelte";
  import { goto } from "$app/navigation";
  import renderMarkdown from "$lib/api/markdown";
  import ReportButton from "./ReportButton.svelte";
  import ProfileLink from "./ProfileLink.svelte";

  export let data: PageData;

  let bio: string | null;
  $: bio = renderMarkdown(data.bio);

  let memberPage: number = 0;
  let memberSlice: PartialMember[] = [];
  $: memberSlice = data.members.slice(memberPage * 20, (memberPage + 1) * 20);
  const totalPages = Math.floor(data.members.length / 20) + 1;

  const prevPage = () => {
    if (memberPage === 0) {
      return;
    }
    memberPage = memberPage - 1;
  };

  const nextPage = () => {
    if ((memberPage + 1) * 20 > data.members.length) {
      return;
    }
    memberPage = memberPage + 1;
  };

  let modalOpen = false;
  let toggleModal = () => (modalOpen = !modalOpen);
  let newMemberName = "";
  let newMemberError: APIError | null = null;

  const createMember = async () => {
    try {
      const member = await apiFetchClient<Member>("/members", "POST", {
        name: newMemberName,
      });

      newMemberName = "";
      newMemberError = null;
      data.members = [...data.members, member];

      toggleModal();

      goto(`/@${data.name}/${member.name}`);
    } catch (e) {
      newMemberError = e as APIError;
    }
  };

  const favNames = data.names.filter((entry) => entry.status === WordStatus.Favourite);
  const favPronouns = data.pronouns.filter((entry) => entry.status === WordStatus.Favourite);

  let profileEmpty = false;
  $: profileEmpty =
    data.names.length === 0 &&
    data.pronouns.length === 0 &&
    data.fields.length === 0 &&
    (!data.bio || data.bio.length === 0);
</script>

<div class="container">
  {#if $userStore && $userStore.id === data.id}
    <Alert color="secondary" fade={false}>
      You are currently viewing your <strong>public</strong> profile.
      <br /><a href="/edit/profile">Edit your profile</a>
    </Alert>
  {/if}
  <div class="grid row-gap-3">
    <div class="row">
      <div class="col-md-4 text-center">
        <FallbackImage width={200} urls={userAvatars(data)} alt="Avatar for @{data.name}" />
      </div>
      <div class="col-md">
        {#if data.display_name}
          <div>
            <h2>{data.display_name}</h2>
            <h5 class="text-body-secondary">@{data.name}</h5>
          </div>
        {:else}
          <h2>@{data.name}</h2>
        {/if}
        {#if profileEmpty && $userStore?.id === data.id}
          <hr />
          <p>
            <em>
              Your profile is empty! You can customize it by going to the <a href="/edit/profile"
                >edit profile</a
              > page.</em
            > <span class="text-muted">(only you can see this)</span>
          </p>
        {:else if bio}
          <hr />
          <p>{@html bio}</p>
        {/if}
      </div>
      {#if data.links.length > 0}
        <div class="col-md d-flex align-items-center">
          <ul class="list-unstyled">
            {#each data.links as link}
              <li><ProfileLink {link} /></li>
            {/each}
          </ul>
        </div>
      {/if}
    </div>
    <div class="row row-cols-1 row-cols-sm-2 row-cols-md-3">
      {#if data.names.length > 0}
        <div class="col-md">
          <h3>Names</h3>
          <ul class="list-unstyled fs-5">
            {#each data.names as name}
              <li><StatusIcon status={name.status} /> {name.value}</li>
            {/each}
          </ul>
        </div>
      {/if}
      {#if data.pronouns.length > 0}
        <div class="col-md">
          <h3>Pronouns</h3>
          <ul class="list-unstyled fs-5">
            {#each data.pronouns as pronouns}
              <li>
                <StatusIcon status={pronouns.status} />
                <PronounLink {pronouns} />
              </li>
            {/each}
          </ul>
        </div>
      {/if}
      {#each data.fields as field}
        <div class="col">
          <FieldCard {field} />
        </div>
      {/each}
    </div>
    {#if $userStore && $userStore.id !== data.id}
      <div class="row">
        <ReportButton subject="user" reportUrl="/users/{data.id}/reports" />
      </div>
    {/if}
    {#if data.members.length > 0 || ($userStore && $userStore.id === data.id)}
      <div class="row">
        <div class="col">
          <hr />
          <h2>
            Members
            {#if $userStore && $userStore.id === data.id}
              <Button
                color="success"
                disabled={data.members.length >= MAX_MEMBERS}
                on:click={toggleModal}><Icon name="person-plus-fill" /> Create member</Button
              >
            {/if}
            {#if totalPages > 1}
              <ButtonGroup>
                <Button on:click={prevPage} disabled={memberPage === 0}
                  ><Icon name="chevron-left" /> Previous page</Button
                >
                <Button disabled>Page {memberPage + 1}/{totalPages}</Button>
                <Button on:click={nextPage} disabled={memberPage === totalPages - 1}
                  >Next page <Icon name="chevron-right" /></Button
                >
              </ButtonGroup>
            {/if}
          </h2>
        </div>
      </div>
      {#if data.members.length > 0}
        <div class="row row-cols-1 row-cols-sm-2 row-cols-md-4 text-center">
          {#each memberSlice as member}
            <PartialMemberCard user={data} {member} />
          {/each}
        </div>
      {:else}
        <div>
          <p>
            You don't have any members yet.
            <br />
            Members are sub-profiles that can have their own avatar,
            names, pronouns, and preferred terms. <span class="text-muted"
              >(only you can see this)</span
            >
          </p>
        </div>
      {/if}
    {/if}
    <Modal header="Create member" isOpen={modalOpen} toggle={toggleModal}>
      <ModalBody>
        <Input bind:value={newMemberName} />
        {#if newMemberError}
          <ErrorAlert error={newMemberError} />
        {/if}
      </ModalBody>
      <ModalFooter>
        <Button color="primary" on:click={createMember} disabled={newMemberName.length === 0}
          >Create member</Button
        >
        <Button color="secondary" on:click={toggleModal}>Cancel</Button>
      </ModalFooter>
    </Modal>
  </div>
</div>

<svelte:head>
  <title>@{data.name} - pronouns.cc</title>

  <meta
    property="og:title"
    content={data.display_name ? `${data.display_name} (@${data.name})` : `@${data.name}`}
  />
  <meta property="og:url" content="{PUBLIC_BASE_URL}/@{data.name}" />

  {#if data.avatar}
    <meta property="og:image" content={userAvatars(data)[0]} />
  {/if}

  {#if favNames.length !== 0 && favPronouns.length !== 0}
    <meta
      property="og:description"
      content="@{data.name} goes by {favNames.map((x) => x.value).join(', ')} and uses {favPronouns
        .map((x) => pronounDisplay(x))
        .join(', ')} pronouns."
    />
  {:else if favNames.length !== 0}
    <meta
      property="og:description"
      content="@{data.name} goes by {favNames.map((x) => x.value).join(', ')}."
    />
  {:else if favPronouns.length !== 0}
    <meta
      property="og:description"
      content="@{data.name} uses {favPronouns.map((x) => pronounDisplay(x)).join(', ')} pronouns."
    />
  {:else if data.bio && data.bio !== ""}
    <meta
      property="og:description"
      content="{data.bio.slice(0, 500)}{data.bio.length > 500 ? '…' : ''}"
    />
  {:else}
    <meta property="og:description" content="@{data.name} on pronouns.cc" />
  {/if}
</svelte:head>
