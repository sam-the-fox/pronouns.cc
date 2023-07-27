<script lang="ts">
  import { flagURL, type PrideFlag } from "$lib/api/entities";
  import { Tooltip } from "sveltestrap";

  export let flag: PrideFlag;

  let elem: HTMLElement;
</script>

<span class="mx-2 my-1">
  <Tooltip target={elem} aria-hidden placement="top">{flag.description ?? flag.name}</Tooltip>
  <img bind:this={elem} class="flag" src={flagURL(flag)} alt={flag.description ?? flag.name} />
  {flag.name}
</span>

<style>
  .flag {
    height: 1.5rem;
    max-width: 200px;
    border-radius: 3px;
  }
</style>
