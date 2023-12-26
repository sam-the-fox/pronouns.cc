<script lang="ts">
  import { Tooltip } from "@sveltestrap/sveltestrap";

  let icon: HTMLElement;
</script>

<Tooltip target={icon} placement="top">This user is an admin</Tooltip>

<div class="profile-badge" bind:this={icon}>
  <svg xmlns="http://www.w3.org/2000/svg" width="1em" height="1em" viewBox="1 1 15.767 14.761"
    ><path
      fill="currentColor"
      d="M3.829.497A2.823 2.823 0 0 0 1 3.326v10.11a2.822 2.822 0 0 0 2.829 2.828h10.11a2.822 2.822 0 0 0 2.829-2.828V3.326a2.823 2.823 0 0 0-2.83-2.83zm6.103 2.427a.33.33 0 0 1 .316.233l1.099 3.603 2.943 1.09a.33.33 0 0 1 0 .62l-2.942 1.088-1.099 3.712a.33.33 0 0 1-.635 0L8.515 9.558 5.573 8.47a.33.33 0 0 1 0-.62l2.943-1.09 1.1-3.601a.33.33 0 0 1 .316-.234zm-5.073.111a.33.33 0 0 1 .31.216l.29.782.782.29a.331.331 0 0 1 0 .62l-.783.29-.289.782a.33.33 0 0 1-.62 0l-.29-.782-.782-.29a.33.33 0 0 1 0-.62l.782-.29.29-.782a.33.33 0 0 1 .31-.216zM6.294 9.87c.149 0 .28.1.319.244l.297 1.091.765.283a.33.33 0 0 1 0 .62l-.76.283-.3 1.197a.331.331 0 0 1-.642 0l-.3-1.197-.762-.282a.33.33 0 0 1 0-.62l.766-.284.297-1.091a.33.33 0 0 1 .32-.244z"
    /></svg
  >
</div>

<style>
  .profile-badge {
    display: inline;
  }
</style>
