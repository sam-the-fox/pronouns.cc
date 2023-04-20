<script lang="ts">
  import { onMount } from "svelte";
  import { browser } from "$app/environment";
  import { decodeJwt } from "jose";

  import {
    Badge,
    Collapse,
    Icon,
    Nav,
    Navbar,
    NavbarBrand,
    NavbarToggler,
    NavItem,
    NavLink,
  } from "sveltestrap";

  import Logo from "./Logo.svelte";
  import { userStore, themeStore, CURRENT_CHANGELOG } from "$lib/store";
  import {
    ErrorCode,
    type APIError,
    type MeUser,
    type Report,
    type Warning,
  } from "$lib/api/entities";
  import { apiFetch, apiFetchClient } from "$lib/api/fetch";
  import { addToast } from "$lib/toast";

  export let commit: string;

  let theme: string;
  let currentUser: MeUser | null;
  let showMenu = false;

  let isAdmin = false;
  let numReports = 0;
  let numWarnings = 0;
  let changelogRead = "99.99.99";

  $: currentUser = $userStore;
  $: theme = $themeStore;

  onMount(() => {
    changelogRead = localStorage.getItem("changelog-read") || "0.0.0";

    const localUser = localStorage.getItem("pronouns-user");
    userStore.set(localUser ? JSON.parse(localUser) : null);

    const token = localStorage.getItem("pronouns-token");
    if (token) {
      apiFetch<MeUser>("/users/@me", { token })
        .then((user) => {
          userStore.set(user);
          localStorage.setItem("pronouns-user", JSON.stringify(user));
        })
        .catch((e) => {
          console.log("getting /users/@me:", e);

          if (
            (e as APIError).code == ErrorCode.InvalidToken ||
            (e as APIError).code == ErrorCode.Forbidden
          ) {
            localStorage.removeItem("pronouns-token");
            localStorage.removeItem("pronouns-user");
          }
        });

      isAdmin = !!decodeJwt(token)["adm"];
      if (isAdmin) {
        apiFetchClient<Report[]>("/admin/reports")
          .then((reports) => {
            numReports = reports.length;
          })
          .catch((e) => {
            console.log("getting reports:", e);
          });
      }

      apiFetchClient<Warning[]>("/auth/warnings")
        .then((warnings) => {
          if (warnings.length !== 0) {
            numWarnings = warnings.length;
            addToast({
              header: "Warnings",
              body: "You have unread warnings. Go to your settings to view them.",
              duration: -1,
            });
          }
        })
        .catch((e) => {
          console.log("getting warnings:", e);
        });
    }
  });

  $: updateTheme(theme);

  const updateTheme = (newTheme: string) => {
    if (!browser) return;

    document.documentElement.setAttribute("data-bs-theme", newTheme);
    localStorage.setItem("pronouns-theme", newTheme);
  };

  const toggleTheme = () => {
    themeStore.set(theme === "dark" ? "light" : "dark");
  };
  const toggleMenu = () => {
    showMenu = !showMenu;
  };
</script>

<Navbar
  color={theme === "dark" ? "dark" : "light"}
  light={theme !== "dark"}
  dark={theme === "dark"}
  expand="lg"
  class="mb-4"
>
  <NavbarBrand href="/" aria-label="pronouns.cc landing page">
    <Logo />
    {#if commit === "[unknown]"}
      <Badge color="danger">DEV</Badge>
    {/if}
  </NavbarBrand>
  <NavbarToggler on:click={toggleMenu} aria-label="Toggle menu" />
  <Collapse isOpen={showMenu} navbar expand="lg">
    <Nav class="ms-auto" navbar>
      {#if currentUser}
        <NavItem>
          <NavLink href="/@{currentUser.name}">@{currentUser.name}</NavLink>
        </NavItem>
        <NavItem>
          <NavLink href="/settings">
            Settings
            {#if numWarnings}
              <Badge color="danger">{numWarnings}</Badge>
            {/if}
          </NavLink>
        </NavItem>
        {#if isAdmin}
          <NavItem>
            <NavLink href="/reports">
              Reports
              {#if numReports !== 0}
                <Badge color="danger">{numReports}</Badge>
              {/if}
            </NavLink>
          </NavItem>
        {/if}
        {#if changelogRead < CURRENT_CHANGELOG}
          <NavItem>
            <NavLink href="/page/changelog">
              Changelog <Badge color="secondary">v{CURRENT_CHANGELOG}</Badge>
            </NavLink>
          </NavItem>
        {/if}
      {:else}
        <NavItem>
          <NavLink href="/auth/login">Log in</NavLink>
        </NavItem>
      {/if}
      <NavItem>
        <NavLink
          on:click={() => toggleTheme()}
          title={theme === "dark" ? "Switch to light mode" : "Switch to dark mode"}
        >
          <Icon name={theme === "dark" ? "sun" : "moon-stars"} height="24" />
          {theme === "dark" ? "Light mode" : "Dark mode"}
        </NavLink>
      </NavItem>
    </Nav>
  </Collapse>
</Navbar>
