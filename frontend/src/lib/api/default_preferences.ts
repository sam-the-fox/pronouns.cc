import { type CustomPreferences, PreferenceSize } from "./entities";

const defaultPreferences: CustomPreferences = {
  favourite: {
    icon: "heart-fill",
    tooltip: "Favourite",
    size: PreferenceSize.Large,
    muted: false,
    favourite: true,
  },
  okay: {
    icon: "hand-thumbs-up",
    tooltip: "Okay",
    size: PreferenceSize.Normal,
    muted: false,
    favourite: false,
  },
  jokingly: {
    icon: "emoji-laughing",
    tooltip: "Jokingly",
    size: PreferenceSize.Normal,
    muted: false,
    favourite: false,
  },
  friends_only: {
    icon: "people",
    tooltip: "Friends only",
    size: PreferenceSize.Normal,
    muted: false,
    favourite: false,
  },
  avoid: {
    icon: "hand-thumbs-down",
    tooltip: "Avoid",
    size: PreferenceSize.Small,
    muted: true,
    favourite: false,
  },
  missing: {
    icon: "question-lg",
    tooltip: "Unknown (missing)",
    size: PreferenceSize.Normal,
    muted: false,
    favourite: false,
  },
};

export default defaultPreferences;
