<script lang="ts">
  import type { PageData } from "./$types";
  import { onMount } from "svelte";

  import {
    Alert,
    Badge,
    Button,
    ButtonGroup,
    Icon,
    Input,
    InputGroup,
    Modal,
    ModalBody,
    ModalFooter,
  } from "sveltestrap";
  import FieldCard from "$lib/components/FieldCard.svelte";
  import PronounLink from "$lib/components/PronounLink.svelte";
  import PartialMemberCard from "$lib/components/PartialMemberCard.svelte";
  import FallbackImage from "$lib/components/FallbackImage.svelte";
  import { userStore } from "$lib/store";
  import {
    MAX_MEMBERS,
    pronounDisplay,
    userAvatars,
    type APIError,
    type Member,
    type PartialMember,
    type CustomPreferences,
    type FieldEntry,
    type Pronoun,
  } from "$lib/api/entities";
  import { PUBLIC_BASE_URL, PUBLIC_SHORT_BASE } from "$env/static/public";
  import { apiFetchClient } from "$lib/api/fetch";
  import ErrorAlert from "$lib/components/ErrorAlert.svelte";
  import { goto } from "$app/navigation";
  import { renderMarkdown } from "$lib/utils";
  import ReportButton from "./ReportButton.svelte";
  import ProfileLink from "./ProfileLink.svelte";
  import { memberNameRegex } from "$lib/api/regex";
  import StatusLine from "$lib/components/StatusLine.svelte";
  import defaultPreferences from "$lib/api/default_preferences";
  import { addToast } from "$lib/toast";
  import ProfileFlag from "./ProfileFlag.svelte";
  import IconButton from "$lib/components/IconButton.svelte";
  import Badges from "./badges/Badges.svelte";

  export let data: PageData;

  let bio: string | null;
  $: bio = renderMarkdown(data.bio);

  let memberPage = 0;
  let memberSlice: PartialMember[] = [];
  $: memberSlice = data.members.slice(memberPage * 20, (memberPage + 1) * 20);
  let totalPages: number;
  $: totalPages = Math.ceil(data.members.length / 20);

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
  let memberNameValid = true;
  $: memberNameValid = memberNameRegex.test(newMemberName);

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

  let mergedPreferences: CustomPreferences;
  $: mergedPreferences = Object.assign({}, defaultPreferences, data.custom_preferences);

  let favNames: FieldEntry[];
  $: favNames = data.names.filter(
    (entry) => (mergedPreferences[entry.status] || { favourite: false }).favourite,
  );
  let favPronouns: Pronoun[];
  $: favPronouns = data.pronouns.filter(
    (entry) => (mergedPreferences[entry.status] || { favourite: false }).favourite,
  );

  let profileEmpty = false;
  $: profileEmpty =
    data.names.length === 0 &&
    data.pronouns.length === 0 &&
    data.fields.length === 0 &&
    (!data.bio || data.bio.length === 0);

  const copyURL = async () => {
    const url = `${PUBLIC_BASE_URL}/@${data.name}`;
    await navigator.clipboard.writeText(url);
    addToast({ body: "Copied the link to your clipboard!", duration: 2000 });
  };

  const copyShortURL = async () => {
    const url = `${PUBLIC_SHORT_BASE}/${data.sid}`;
    await navigator.clipboard.writeText(url);
    addToast({ body: "Copied the short link to your clipboard!", duration: 2000 });
  };

  onMount(async () => {
    if ($userStore && $userStore.id === data.id) {
      console.log("User is current user, fetching members");
      try {
        const members = await apiFetchClient<PartialMember[]>("/users/@me/members");
        data.members = members;
      } catch (e) {
        // If it fails, we fail silently but log to console anyway
        console.error("Fetching members:", e);
      }
    }
  });
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
        {#if data.flags && data.bio}
          <div class="d-flex flex-wrap m-4">
            {#each data.flags as flag}
              <ProfileFlag {flag} />
            {/each}
          </div>
        {/if}
      </div>
      <div class="col-md">
        {#if data.display_name}
          <div>
            <h2>{data.display_name} <Badges userBadges={data.badges} /></h2>
            <p class="fs-5 text-body-secondary">@{data.name}</p>
          </div>
        {:else}
          <h2>@{data.name} <Badges userBadges={data.badges} /></h2>
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
              <ProfileLink {link} />
            {/each}
          </ul>
        </div>
      {/if}
    </div>
    {#if data.flags && !data.bio}
      <div class="d-flex flex-wrap m-4">
        {#each data.flags as flag}
          <ProfileFlag {flag} />
        {/each}
      </div>
    {/if}
    <div class="row row-cols-1 row-cols-sm-2 row-cols-md-3">
      {#if data.names.length > 0}
        <div class="col-md">
          <h3>Names</h3>
          <ul class="list-unstyled fs-5">
            {#each data.names as name}
              <li>
                <StatusLine preferences={data.custom_preferences} status={name.status}
                  >{name.value}</StatusLine
                >
              </li>
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
                <StatusLine preferences={data.custom_preferences} status={pronouns.status}
                  ><PronounLink {pronouns} /></StatusLine
                >
              </li>
            {/each}
          </ul>
        </div>
      {/if}
      {#each data.fields as field}
        <div class="col">
          <FieldCard preferences={data.custom_preferences} {field} />
        </div>
      {/each}
    </div>
    <div class="row">
      <div class="col-md-6">
        <InputGroup>
          <Button color="secondary" outline on:click={copyURL}>
            <Icon name="clipboard" /> Copy link
          </Button>
          {#if PUBLIC_SHORT_BASE}
            <IconButton
              outline
              icon="link-45deg"
              tooltip="Copy short link"
              color="secondary"
              click={copyShortURL}
            />
          {/if}
          {#if $userStore && $userStore.id !== data.id}
            <ReportButton subject="user" reportUrl="/users/{data.id}/reports" />
          {/if}
        </InputGroup>
      </div>
      <div class="col-md-6" />
    </div>
    {#if data.members.length > 0 || ($userStore && $userStore.id === data.id)}
      <div class="row">
        <div class="col">
          <hr />
          <h2>
            {data.member_title || "Members"}
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
        <div class="row row-cols-2 row-cols-md-3 row-cols-lg-4 row-cols-xl-5 text-center">
          {#each memberSlice as member}
            <PartialMemberCard user={data} {member} />
          {/each}
        </div>
        {#if totalPages > 1}
          <div class="text-center">
            <ButtonGroup>
              <Button on:click={prevPage} disabled={memberPage === 0}
                ><Icon name="chevron-left" /> Previous page</Button
              >
              <Button disabled>Page {memberPage + 1}/{totalPages}</Button>
              <Button on:click={nextPage} disabled={memberPage === totalPages - 1}
                >Next page <Icon name="chevron-right" /></Button
              >
            </ButtonGroup>
          </div>
        {/if}
      {:else}
        <div>
          <p>
            You don't have any members yet.
            <br />
            Members are sub-profiles that can have their own avatar, names, pronouns, and preferred terms.
            <br />
            If you were expecting to see members here, check your
            <a href="/settings/members">list of hidden members</a>.
            <span class="text-muted">(only you can see this)</span>
          </p>
        </div>
      {/if}
    {/if}
    <Modal header="Create member" isOpen={modalOpen} toggle={toggleModal}>
      <ModalBody>
        <Input bind:value={newMemberName} />
        <p class="text-muted my-2">
          <Icon name="info-circle-fill" aria-label="Info" /> Your members must have distinct names. Member
          names must be 100 characters long at most, and cannot contain the following characters: @ ?
          ! # / \ [ ] " ' $ % & ( ) &lbrace; &rbrace; + &lt; = &gt; ^ | ~ ` and ,
        </p>
        {#if newMemberError}
          <ErrorAlert error={newMemberError} />
        {/if}
      </ModalBody>
      <ModalFooter>
        {#if !memberNameValid && newMemberName.length > 0}
          <span class="text-danger-emphasis mb-2">That member name is not valid.</span>
        {/if}
        <Button color="primary" on:click={createMember} disabled={!memberNameValid}
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
  <meta name="description" content="@{data.name} on pronouns.cc" />
  <meta name="robots" content="noindex" />

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
