<script lang="ts">
  import { flagURL, type PrideFlag } from "$lib/api/entities";
  import { Button, Tooltip } from "sveltestrap";

  export let flag: PrideFlag;
  export let tooltip: string;
  let className: string | null | undefined = undefined;
  export { className as class };

  let elem: HTMLElement;
</script>

<Tooltip target={elem} placement="top">{tooltip}</Tooltip>
<Button bind:inner={elem} class={className} on:click color="secondary" outline>
  <img class="flag" src={flagURL(flag)} alt={flag.description ?? flag.name} />
  {flag.name}
</Button>

<style>
  .flag {
    height: 1.5rem;
    max-width: 200px;
    border-radius: 3px;
    margin-left: -5px;
  }
</style>
