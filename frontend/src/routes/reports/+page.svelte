<script lang="ts">
  import { DateTime } from "luxon";
  import { Button, Card, CardBody, CardFooter, CardHeader } from "sveltestrap";
  import type { PageData } from "./$types";

  export let data: PageData;
</script>

<svelte:head>
  <title>Reports - pronouns.cc</title>
</svelte:head>

<div class="container">
  <h1>Reports</h1>

  <div>
    {#each data.reports as report}
      <Card>
        <CardHeader>
          <strong>#{report.id}</strong> on <a href="/@{report.user_name}">@{report.user_name}</a>
          ({report.user_id}) {#if report.member_id}
            (member: <a href="/@{report.user_name}/{report.member_name}">{report.member_name}</a>,
            {report.member_id})
          {/if}
        </CardHeader>
        <CardBody>
          <blockquote class="blockquote">{report.reason}</blockquote>
        </CardBody>
        <CardFooter>
          Created {DateTime.fromISO(report.created_at)
            .toLocal()
            .toLocaleString(DateTime.DATETIME_MED)} &bull;
          <Button outline color="warning" size="sm">Warn user</Button>
          <Button outline color="danger" size="sm">Deactivate user</Button>
          <Button outline color="secondary" size="sm">Ignore report</Button>
        </CardFooter>
      </Card>
    {/each}
  </div>
</div>
