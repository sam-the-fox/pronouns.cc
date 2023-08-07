<script lang="ts">
  import { NavLink } from "sveltestrap";
  import { page } from "$app/stores";

  export let href: string;
</script>

<NavLink {href} active={$page.url.pathname === href}>
  <slot />
</NavLink>
