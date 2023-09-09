<script lang="ts">
  import { page } from "$app/stores";
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
  import { userStore, themeStore, CURRENT_CHANGELOG, settingsStore } from "$lib/store";
  import {
    ErrorCode,
    type APIError,
    type MeUser,
    type Report,
    type Warning,
    type Settings,
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

  $: currentUser = $userStore;
  $: theme = $themeStore;

  onMount(() => {
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

      apiFetchClient<Settings>("/users/@me/settings", "GET", null, 2)
        .then((data) => {
          settingsStore.set({ current: true, settings: data });
        })
        .catch((e) => {
          console.log("getting user settings:", e);
        });

      // TODO: is there a cleaner way to do this? also, remove this eventually
      const oldChangelogRead = localStorage.getItem("changelog-read");
      const oldSettingsNotice = localStorage.getItem("alert-1681976313");
      if (oldChangelogRead || oldSettingsNotice) {
        localStorage.removeItem("changelog-read");
        localStorage.removeItem("alert-1681976313");

        apiFetchClient<Settings>(
          "/users/@me/settings",
          "PATCH",
          {
            read_changelog: oldChangelogRead ? oldChangelogRead : undefined,
            read_settings_notice: oldSettingsNotice ? "1681976313" : undefined,
          },
          2,
        )
          .then((data) => {
            settingsStore.set({ current: true, settings: data });
          })
          .catch((e) => {
            console.log("updating user settings:", e);
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
      <strong id="beta-text" class="text-danger">dev</strong>
    {:else}
      <span id="beta-text">beta</span>
    {/if}
  </NavbarBrand>
  <NavbarToggler on:click={toggleMenu} aria-label="Toggle menu" />
  <Collapse isOpen={showMenu} navbar expand="lg">
    <Nav class="ms-auto" navbar>
      {#if currentUser}
        <NavItem>
          <NavLink
            href="/@{currentUser.name}"
            active={$page.url.pathname === `/@${currentUser.name}`}
          >
            @{currentUser.name}
          </NavLink>
        </NavItem>
        <NavItem>
          <NavLink href="/settings" active={$page.url.pathname.startsWith("/settings")}>
            Settings
            {#if numWarnings}
              <Badge color="danger">{numWarnings}</Badge>
            {/if}
          </NavLink>
        </NavItem>
        {#if isAdmin}
          <NavItem>
            <NavLink href="/reports" active={$page.url.pathname.startsWith("/reports")}>
              Reports
              {#if numReports !== 0}
                <Badge color="danger">{numReports}</Badge>
              {/if}
            </NavLink>
          </NavItem>
        {/if}
        {#if $settingsStore.current && $settingsStore.settings.read_changelog < CURRENT_CHANGELOG}
          <NavItem>
            <NavLink href="/page/changelog" active={$page.url.pathname === "/page/changelog"}>
              Changelog <Badge color="secondary">v{CURRENT_CHANGELOG}</Badge>
            </NavLink>
          </NavItem>
        {/if}
      {:else}
        <NavItem>
          <NavLink href="/auth/login" active={$page.url.pathname === "/auth/login"}>Log in</NavLink>
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

<style>
  /* These exact values make it look almost identical to the SVG version, which is what we want */
  #beta-text {
    font-size: 0.7em;
    position: relative;
    font-style: italic;
    bottom: 12px;
    right: 3px;
  }
</style>
