<script lang="ts">
  import { onMount } from "svelte";
  import { browser } from "$app/environment";

  import { Icon } from "@steeze-ui/svelte-icon";
  import { Moon, Sun, EllipsisVertical } from "@steeze-ui/heroicons";

  import NavItem from "./NavItem.svelte";
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
      document.documentElement.classList.remove("dark");
    } else {
      document.documentElement.classList.add("dark");
    }
    localStorage.setItem("theme", isDark ? "dark" : "light");
  };

  const toggleTheme = () => {
    darkTheme = !darkTheme;
  };
</script>

<div
  class="bg-white/75 dark:bg-slate-800/75 w-full backdrop-blur border-slate-200 dark:border-slate-700 border-b"
>
  <div class="max-w-8xl mx-auto">
    <div class="py-4 mx-4">
      <div class="flex items-center">
        <a href="/"><Logo /></a>
        <div class="ml-auto flex items-center">
          <nav class="hidden lg:flex">
            <ul class="flex space-x-4 font-bold">
              <NavItem href="/login">Log in</NavItem>
            </ul>
          </nav>
          <div
            class="flex border-l border-slate-200 ml-4 pl-4 lg:ml-6 lg:pl-6 lg:mr-2 dark:border-slate-700 space-x-2 lg:space-x-4"
          >
            <div
              class="cursor-pointer flex"
              on:click={() => toggleTheme()}
              on:keypress={() => toggleTheme()}
              title={darkTheme ? "Switch to light mode" : "Switch to dark mode"}
            >
              <Icon src={darkTheme ? Sun : Moon} class="hover:text-sky-400" size="32" />
            </div>
            <div
              class="cursor-pointer flex lg:hidden"
              on:click={() => (showMenu = !showMenu)}
              on:keypress={() => (showMenu = !showMenu)}
              title={showMenu ? "Hide menu" : "Show menu"}
            >
              <Icon src={EllipsisVertical} class="hover:text-sky-400" size="32" />
            </div>
            <!-- <div
            onClick={() => setShowMenu(!showMenu)}
            title="Show menu"
            className="cursor-pointer flex lg:hidden"
          >
            <List
              className="dark:hover:text-sky-400 hover:text-sky-500"
              size={24}
            />
          </div> -->
          </div>
        </div>
      </div>
    </div>
  </div>
</div>
