import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";

interface Pronouns {
  [key: string]: { pronouns: string[]; display?: string };
}

export const load = (async ({ params }) => {
  const [param, displayText] = params.pronouns.split(",");

  const arr = param.split("/");
  if (arr.length === 0 || params.pronouns === "") {
    throw error(404, "Pronouns not found");
  }

  if (arr.length === 5) {
    return {
      displayText,
      subjective: arr[0],
      objective: arr[1],
      possessiveDeterminer: arr[2],
      possessivePronoun: arr[3],
      reflexive: arr[4],
    };
  }

  const pronouns: Pronouns = (await import("$lib/pronouns.json")).pronouns;

  if (params.pronouns in pronouns) {
    const entry = pronouns[params.pronouns];
    return {
      displayText: entry.display || "",
      subjective: entry.pronouns[0],
      objective: entry.pronouns[1],
      possessiveDeterminer: entry.pronouns[2],
      possessivePronoun: entry.pronouns[3],
      reflexive: entry.pronouns[4],
    };
  }

  throw error(404, "Pronouns not found");
}) satisfies PageLoad;
