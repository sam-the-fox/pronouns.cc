import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";
import { ErrorCode, type PronounsJson } from "$lib/api/entities";

import pronounsRaw from "$lib/pronouns.json";
const pronouns = pronounsRaw as PronounsJson;

export const load = (async ({ params }) => {
  const [param, displayText] = decodeURIComponent(params.pronouns).split(",");

  const arr = param.split("/");
  if (arr.length === 0 || params.pronouns === "") {
    throw error(404, { code: ErrorCode.NotFound, message: "Pronouns not found" });
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

  if (params.pronouns in pronouns.pages) {
    const entry = pronouns.pages[params.pronouns];
    return {
      displayText: entry.display || "",
      subjective: entry.pronouns[0],
      objective: entry.pronouns[1],
      possessiveDeterminer: entry.pronouns[2],
      possessivePronoun: entry.pronouns[3],
      reflexive: entry.pronouns[4],
    };
  } else if (params.pronouns in pronouns.autocomplete) {
    const entry = pronouns.autocomplete[params.pronouns];
    return {
      displayText: entry.display || "",
      subjective: entry.pronouns[0],
      objective: entry.pronouns[1],
      possessiveDeterminer: entry.pronouns[2],
      possessivePronoun: entry.pronouns[3],
      reflexive: entry.pronouns[4],
    };
  }

  throw error(404, { code: ErrorCode.NotFound, message: "Pronouns not found" });
}) satisfies PageLoad;
