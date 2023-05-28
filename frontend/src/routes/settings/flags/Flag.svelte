<script lang="ts">
  import { flagURL, type PrideFlag } from "$lib/api/entities";
  import { Button } from "sveltestrap";

  export let flag: PrideFlag;
</script>

<Button outline class="m-1">
  <img class="flag" src={flagURL(flag)} alt={flag.description ?? flag.name} />
  {flag.name}
</Button>

<style>
  .flag {
    height: 2rem;
    max-width: 200px;
    border-radius: 3px;
    margin-left: -5px;
  }
</style>
