<script lang="ts">
  import type { Report } from "$lib/api/entities";
  import { DateTime } from "luxon";
  import { Button, Card, CardBody, CardFooter, CardHeader } from "sveltestrap";

  export let report: Report;
</script>

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
    Created {DateTime.fromISO(report.created_at).toLocal().toLocaleString(DateTime.DATETIME_MED)}
    <slot />
  </CardFooter>
</Card>
