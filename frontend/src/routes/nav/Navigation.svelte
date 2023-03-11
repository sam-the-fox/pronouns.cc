<script lang="ts">
  import { onMount } from "svelte";
  import { browser } from "$app/environment";

  import {
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
  import { userStore, themeStore } from "$lib/store";
  import { ErrorCode, type APIError, type MeUser } from "$lib/api/entities";
  import { apiFetch } from "$lib/api/fetch";

  let theme: string;
  let currentUser: MeUser | null;
  let showMenu: boolean = false;

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
    }
  });

  $: updateTheme(theme);

  const updateTheme = (newTheme: string) => {
    if (!browser) return;

    document.documentElement.setAttribute("data-bs-theme", newTheme);
    localStorage.setItem("pronouns-theme", newTheme);
  };

  const toggleTheme = () => {
    themeStore.set(theme === "dark" ? "light" : "dark")
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
  <NavbarBrand href="/"><Logo /></NavbarBrand>
  <NavbarToggler on:click={toggleMenu} />
  <Collapse isOpen={showMenu} navbar expand="lg">
    <Nav class="ms-auto" navbar>
      <NavItem>
        {#if currentUser}
          <NavLink href="/@{currentUser.name}">@{currentUser.name}</NavLink>
          <NavLink href="/settings">Settings</NavLink>
          <NavLink href="/logout">Log out</NavLink>
        {:else}
          <NavLink href="/login">Log in</NavLink>
        {/if}
      </NavItem>
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
