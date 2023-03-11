<script lang="ts">
  import type { PartialMember, User } from "$lib/api/entities";
  import FallbackImage from "./FallbackImage.svelte";

  export let user: User;
  export let member: PartialMember;
</script>

<div>
  <FallbackImage
    urls={member.avatar_urls}
    alt="Avatar for {member.name}"
  />
  <a class="text-reset" href="/@{user.name}/{member.name}"><h5 class="m-2">{member.display_name ?? member.name}</h5></a>
</div>
