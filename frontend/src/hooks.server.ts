import { PRIVATE_SENTRY_DSN } from "$env/static/private";
import * as Sentry from "@sentry/node";
import type { HandleServerError } from "@sveltejs/kit";
import { randomUUID } from "node:crypto";

if (PRIVATE_SENTRY_DSN) {
  Sentry.init({ dsn: PRIVATE_SENTRY_DSN });
}

export const handleError = (({ error }) => {
  let status = 500;
  let id = randomUUID().replaceAll("-", "");
  if (typeof error === "object" && error && "status" in error) {
    if (error.status === 500 && PRIVATE_SENTRY_DSN) {
      id = Sentry.captureException(error);
    } else {
      status = error.status as number;
    }
  }

  return {
    message: "Internal server error",
    code: status,
    id,
  };
}) satisfies HandleServerError;
