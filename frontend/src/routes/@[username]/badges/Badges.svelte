<script lang="ts">
  import Admin from "./Admin.svelte";

  export let userBadges: number;

  let isAdmin: boolean;
  $: isAdmin = (userBadges & (1 << 0)) === 1 << 0;
</script>

<div class="profile-badges">
  {#if isAdmin}
    <Admin />
  {/if}
</div>

<style>
  .profile-badges {
    display: inline;
  }
</style>
