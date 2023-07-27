// See https://kit.svelte.dev/docs/types#app
// for information about these interfaces
declare global {
  namespace App {
    // interface Error {}
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
