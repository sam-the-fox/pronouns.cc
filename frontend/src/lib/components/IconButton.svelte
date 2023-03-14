<script lang="ts">
  import { Button, Icon, Tooltip } from "sveltestrap";

  export let icon: string;
  export let color: "primary" | "secondary" | "success" | "danger";
  export let tooltip: string;
  export let active: boolean = false;
  export let disabled: boolean = false;

  export let click: ((e: MouseEvent) => void) | undefined = undefined;
  export let href: string | undefined = undefined;

  let button: HTMLElement;
</script>

<Tooltip target={button} placement="top">{tooltip}</Tooltip>
<Button {color} {active} {disabled} {href} on:click={click} bind:inner={button}>
  <Icon name={icon} />
</Button>
