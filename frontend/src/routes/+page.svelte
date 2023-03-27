<script lang="ts">
  import { PUBLIC_BASE_URL } from "$env/static/public";
  import { Button } from "sveltestrap";
  import { userStore } from "$lib/store";
</script>

<svelte:head>
  <title>pronouns.cc</title>
  <meta name="description" content="Create and share your preferred names and pronouns!" />
  <meta property="og:title" content="pronouns.cc" />
  <meta property="og:description" content="Create and share your preferred names and pronouns!" />
  <meta property="og:url" content={PUBLIC_BASE_URL} />
</svelte:head>

<div>
  <h1><b>pronouns</b>.cc</h1>

  <p class="fs-3">Create and share <em>all</em> your pronouns</p>

  <p>
    {#if $userStore}
      <Button href="/@{$userStore.name}" color="secondary" size="lg">Go to your profile</Button>
    {:else}
      <Button href="/auth/login" color="success" size="lg">Log in or sign up</Button>
    {/if}
  </p>
</div>

<div class="row">
  <div class="col fs-5">
    <p>
      Everybody uses pronouns to refer to others. Whether it be he/him, she/her, they/them, or
      anything else (or nothing at all!), everybody has preferences. But how do you know which
      pronouns you should use?
    </p>
    <p>
      pronouns.cc is a service where you can create and share a list of your preferred names,
      pronouns, and terms to share with other people! But not just <em>one</em> set: you can create
      <em>multiple</em>
      lists and have all of them linked on your main profile.
    </p>
    <p>This is useful for plural systems, people who kin, and anyone else too!</p>
  </div>
</div>

<div class="row row-cols-1 row-cols-md-3">
  <div>
    <h5>In beta</h5>
    <p>
      pronouns.cc is currently <strong>in beta</strong>. There might be issues and some
      functionality is not available or unfinished. Issue reports and pull requests
      <a href="https://codeberg.org/u1f320/pronouns.cc" rel="noreferrer" target="_blank"
        >in the repository</a
      > are welcome!
    </p>
  </div>
  <div>
    <h5>Isn't this similar to Pronouns.page?</h5>
    <p>
      Yes! <a href="https://pronouns.page" rel="noreferrer" target="_blank">Pronouns.page</a>, along
      with <a href="https://pronouny.xyz" rel="noreferrer" target="_blank">pronouny.xyz</a>, was a
      direct inspiration for pronouns.cc. If those sites are working fine for you, that's great,
      they're good projects!
    </p>
  </div>
  <div>
    <h5>Open source</h5>
    <p>
      pronouns.cc is
      <a href="https://codeberg.org/u1f320/pronouns.cc" rel="noreferrer" target="_blank"
        >open source</a
      >, and licensed under the GNU Affero General Public License. Feel free to contribute!
    </p>
  </div>
  <div>
    <h5>This is useful!</h5>
    <p>thank you :3</p>
  </div>
  <div>
    <h5>This website sucks!</h5>
    <p>oh no :&lpar;</p>
  </div>
  <div>
    <h5>Support pronouns.cc</h5>
    <p>
      If you like pronouns.cc and want to support me financially, check out the <a
        href="/page/about">about page</a
      >!
    </p>
  </div>
</div>
