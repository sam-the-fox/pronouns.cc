<script lang="ts">
  import { Button, Icon, Tooltip } from "sveltestrap";

  export let icon: string;
  export let color: "primary" | "secondary" | "success" | "danger";
  export let tooltip: string;
  export let active: boolean = false;
  export let click: (e: MouseEvent) => void;

  let button: HTMLElement;
</script>

<Tooltip target={button} placement="top">{tooltip}</Tooltip>
<Button {color} {active} on:click={click} bind:inner={button}>
  <Icon name={icon} />
</Button>
