<script lang="ts">
  import { Toast } from "@sveltestrap/sveltestrap";

  export let header: string | undefined = undefined;
  export let body: string;

  let isOpen = true;
  const toggle = () => (isOpen = !isOpen);
</script>

<Toast {isOpen} {header} body class="m-3" {toggle}>
  {body}
</Toast>
