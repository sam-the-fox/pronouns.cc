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

declare module "svelte-hcaptcha" {
  import type { SvelteComponent } from "svelte";

  export interface HCaptchaProps {
    sitekey?: string;
    apihost?: string;
    hl?: string;
    reCaptchaCompat?: boolean;
    theme?: CaptchaTheme;
    size?: string;
  }

  declare class HCaptcha extends SvelteComponent {
    $$prop_def: HCaptchaProps;
  }

  export default HCaptcha;
}

export {};
