<script lang="ts">
  import { goto } from "$app/navigation";
  import type { APIError, MeUser } from "$lib/api/entities";
  import { apiFetch, apiFetchClient } from "$lib/api/fetch";
  import { userStore } from "$lib/store";
  import type { PageData } from "./$types";
  import { addToast } from "$lib/toast";
  import CallbackPage from "../CallbackPage.svelte";
  import type { SignupResponse } from "$lib/api/responses";

  export let data: PageData;

  const signupForm = async (username: string, invite: string) => {
    try {
      const resp = await apiFetch<SignupResponse>("/auth/tumblr/signup", {
        method: "POST",
        body: {
          ticket: data.ticket,
          username: username,
          invite_code: invite,
        },
      });

      localStorage.setItem("pronouns-token", resp.token);
      localStorage.setItem("pronouns-user", JSON.stringify(resp.user));
      userStore.set(resp.user);
      addToast({ header: "Welcome!", body: "Signed up successfully!" });
      goto(`/@${resp.user.name}`);
    } catch (e) {
      data.error = e as APIError;
    }
  };

  const linkAccount = async () => {
    try {
      const resp = await apiFetchClient<MeUser>("/auth/tumblr/add-provider", "POST", {
        ticket: data.ticket,
      });

      localStorage.setItem("pronouns-user", JSON.stringify(resp));
      userStore.set(resp);
      addToast({ header: "Linked account", body: "Successfully linked account!" });
      await goto("/settings/auth");
    } catch (e) {
      data.error = e as APIError;
    }
  };
</script>

<CallbackPage
  authType="Tumblr"
  remoteName={data.tumblr}
  error={data.error}
  requireInvite={data.require_invite}
  isDeleted={data.is_deleted}
  ticket={data.ticket}
  token={data.token}
  user={data.user}
  deletedAt={data.deleted_at}
  selfDelete={data.self_delete}
  deleteReason={data.delete_reason}
  {linkAccount}
  {signupForm}
/>
