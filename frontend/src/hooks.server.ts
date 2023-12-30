import { PRIVATE_SENTRY_DSN } from "$env/static/private";
import * as Sentry from "@sentry/node";
import type { HandleServerError } from "@sveltejs/kit";

if (PRIVATE_SENTRY_DSN) {
  Sentry.init({ dsn: PRIVATE_SENTRY_DSN });
}

export const handleError = (({ error, event }) => {
  console.log(error);
  console.log(event);
  // const id = Sentry.captureException(error);

  return {
    message: "Internal server error",
    code: 500,
    //id,
  };
}) satisfies HandleServerError;
