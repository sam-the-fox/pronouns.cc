// See https://kit.svelte.dev/docs/types#app

import type { ErrorCode } from "$lib/api/entities";

// for information about these interfaces
declare global {
  namespace App {
    interface Error {
      code: ErrorCode;
      message?: string | undefined;
      details?: string | undefined;
    }
    // interface Locals {}
    // interface PageData {}
    // interface Platform {}
  }
}

export {};
