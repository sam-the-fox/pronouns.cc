# Changelog

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
