<script lang="ts">
  import type { PageData } from "../../$types";

  export let data: PageData;

  const { subjective, objective, possessiveDeterminer, possessivePronoun, reflexive } = data;

  const displayText = data.displayText || `${subjective}/${objective}`;
</script>

<h1>{displayText}</h1>

<p>Here are some example sentences using <b>{displayText}</b> pronouns!</p>

<blockquote class="blockquote">
  <p><b class="text-capitalize">{subjective}</b> went to the park.</p>
  <p>I went with <b>{objective}</b>.</p>
  <p><b class="text-capitalize">{subjective}</b> brought <b>{possessiveDeterminer}</b> frisbee.</p>
  <p>At least, I think it was <b>{possessivePronoun}</b>.</p>
  <p>
    <b class="text-capitalize">{objective}</b> threw the frisbee to
    <b>{reflexive}</b>.
  </p>
</blockquote>
