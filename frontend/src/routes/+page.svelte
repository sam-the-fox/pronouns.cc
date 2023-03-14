<script lang="ts">
  import { PUBLIC_BASE_URL } from "$env/static/public";
</script>

<svelte:head>
  <title>pronouns.cc</title>
  <meta property="og:title" content="pronouns.cc" />
  <meta property="og:description" content="Name and pronoun cards!" />
  <meta property="og:url" content={PUBLIC_BASE_URL} />
</svelte:head>

<h1>pronouns.cc</h1>

<p>
  pronouns.cc is a service to create and share pronoun cards, similar to
  <a href="https://pronouns.page" rel="noreferrer" target="_blank">pronouns.page</a> and
  <a href="https://pronouny.xyz" rel="noreferrer" target="_blank">pronouny.xyz</a>.
</p>

<p>
  You're be able to create multiple profiles that are linked together, useful for plurality (<a
    href="https://morethanone.info"
    rel="noreferrer"
    target="_blank">what?</a
  >) and kin, or even just for fun!
</p>

<p>
  pronouns.cc is currently in a <b>closed beta</b>. Sign ups will be opened Soon™. You can check out
  the source code on
  <a href="https://codeberg.org/u1f320/pronouns.cc" rel="noreferrer" target="_blank">Codeberg</a>!
  &lt;3
</p>
