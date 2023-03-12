<script lang="ts">
  import { goto } from "$app/navigation";
  import { userStore } from "$lib/store";
  import { onMount } from "svelte";

  export const ssr = false;

  onMount(() => {
    userStore.set(null);
    localStorage.removeItem("pronouns-token");
    localStorage.removeItem("pronouns-user");
    goto("/");
  });
</script>

<h1>Log out</h1>
