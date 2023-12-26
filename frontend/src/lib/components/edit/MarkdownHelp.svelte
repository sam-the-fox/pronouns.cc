<script lang="ts">
  import { Icon, Modal } from "@sveltestrap/sveltestrap";

  let isOpen = false;
  const toggle = () => (isOpen = !isOpen);
</script>

<Icon name="info-circle-fill" aria-hidden /> Your bio supports limited
<!-- svelte-ignore a11y-invalid-attribute -->
<a class="text-reset" href="" on:click={toggle}>Markdown</a>.

<Modal header="Markdown help" body {isOpen} {toggle}>
  <p>
    Markdown is a simple way to format text, using symbols like <code>*</code>, <code>_</code>, and
    <code>`</code>. If you have used Discord before, formatting works the same as there, just with a
    couple different options.
  </p>
  <hr />
  <p>Type <code>*Italic*</code> or <code>_Italic_</code> to get <em>Italic</em></p>
  <hr />
  <p>Type <code>**Bold**</code> or <code>__Bold__</code> to get <strong>Bold</strong></p>
  <hr />
  <p>Use at least three dashes <code>---</code> on a separate line to add a horizontal rule:</p>
  <hr class="m-3" />
  <hr />
  <p>Type <code>`Inline code`</code> to get <code>Inline code</code></p>
  <p>Use three backticks to make a code block:</p>
  <pre><code>```<br />print("Code here")<br />```</code></pre>
  <hr />
  <p>Type</p>
  <pre><code>* List<br />* List<br />* List</code></pre>
  <p>to get</p>
  <ul>
    <li>List</li>
    <li>List</li>
    <li>List</li>
  </ul>
  <hr />
  <p>Type</p>
  <pre><code>1. List<br />2. List<br />3. List</code></pre>
  <p>to get</p>
  <ol>
    <li>List</li>
    <li>List</li>
    <li>List</li>
  </ol>
</Modal>
