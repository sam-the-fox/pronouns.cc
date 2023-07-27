# Changelog

## 0.5.3 - June 3rd 2023

- Added a link shortener at **prns.cc**! There is now an additional button on all user and member pages to copy a shorter link.
  The IDs used here can be rerolled on the edit profile page.
- Fixed the member list breaking if you have too many private members.

## 0.5.0 - May 29th 2023

- **Added pride flags!**
  These will show up below avatars and are completely customizable:
  you can upload up to one hundred different flags and pick which ones will show up on which profiles.
- `.` and `..` are no longer valid member names, as these break member pages. Any other number of periods is still allowed.
- Added Plausible analytics.
- Added some extra info to the [settings page](/settings).
- Added [API documentation](/api).

## 0.4.2 - May 11th 2023

- The report button will no longer show up for your own members.
- Added a warning on the edit member page if your member list is hidden.
- Fixed hidden members not showing up when viewing your _own_ user page.
- The site will now default to a dark theme while loading, to remove the white flash when using a dark theme.
  Sadly, it's one or the other, so light theme users will now see a _dark_ flash when the site is loading—
  having asked multiple users, the general consensus is that a dark flash is better than a light flash.

## 0.4.1 - April 24th 2023

- Added buttons to change the order of links on profiles.
- Added a "copy link" button to user and member profiles.
- Added a second set of page buttons below the member list.
- Added a captcha when signing up to prevent spam bots.
- Custom preferences with "treat as favourite" checked are now shown in the member list.
- Changed how the member list is displayed: the site will now show between two and five members per row depending on screen size.
- Fixed the save button sometimes not showing up when editing custom preferences.
- Fixed custom preferences not showing up if they were added in the current editing session.

## 0.4.0 - April 20th 2023

- **Added custom preferences!**
  These can be used instead of, or in addition to, the existing five preferences
  and can use any icon in [Bootstrap Icons](https://icons.getbootstrap.com/).
  Change them in your [profile settings](/edit/profile).
- Added the ability to sign in with Tumblr and Google accounts,
  in addition to Discord and Mastodon.
- Added this changelog.
- Added a donation link to the footer, and a message on the settings page. Sadly, running a website like this costs money :(
- The website will now correctly handle tokens expiring, so you don't get stuck in an infinite loop as soon as they expire anymore.
- Fixed pronoun links, pronouns with special characters will link to the correct page now.
