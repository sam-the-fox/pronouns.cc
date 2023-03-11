<script lang="ts">
  import { onMount } from "svelte";
  import { browser } from "$app/environment";

  import { Collapse,Icon,  Nav, Navbar, NavbarBrand, NavbarToggler, NavItem, NavLink } from "sveltestrap";

  import Logo from "./Logo.svelte";

  let darkTheme: boolean = false;
  let showMenu: boolean = false;

  onMount(() => {
    darkTheme =
      localStorage.theme === "dark" ||
      (!("theme" in localStorage) && window.matchMedia("(prefers-color-scheme: dark)").matches);
  });

  $: updateTheme(darkTheme);

  const updateTheme = (isDark: boolean) => {
    if (!browser) return;

    if (isDark) {
      document.documentElement.setAttribute("data-bs-theme", "dark");
    } else {
      document.documentElement.setAttribute("data-bs-theme", "light");
    }
    localStorage.setItem("theme", isDark ? "dark" : "light");
  };

  const toggleTheme = () => {
    darkTheme = !darkTheme;
  };
  const toggleMenu = () => {
    showMenu = !showMenu;
  };
</script>

<Navbar
  color={darkTheme ? "dark" : "light"}
  light={!darkTheme}
  dark={darkTheme}
  expand="lg"
  class="mb-4"
>
  <NavbarBrand href="/"><Logo /></NavbarBrand>
  <NavbarToggler on:click={toggleMenu} />
  <Collapse isOpen={showMenu} navbar expand="lg">
    <Nav class="ms-auto" navbar>
      <NavItem>
        <NavLink href="/login">Log in</NavLink>
      </NavItem>
      <NavItem>
        <NavLink
          on:click={() => toggleTheme()}
          title={darkTheme ? "Switch to light mode" : "Switch to dark mode"}
        >
          <Icon name={darkTheme ? "sun" : "moon-stars"} height="24" />
          {darkTheme ? "Light mode" : "Dark mode"}
        </NavLink>
      </NavItem>
    </Nav>
  </Collapse>
</Navbar>
