import { useEffect, useState } from "react";
import { List, MoonStars, Sun } from "react-bootstrap-icons";

import Link from "next/link";
import { useRecoilState } from "recoil";

import { APIError, ErrorCode, fetchAPI, MeUser } from "../lib/api-fetch";
import { themeState, userState } from "../lib/state";
import Logo from "./logo";
import NavItem from "./NavItem";

export default function Navigation() {
  const [user, setUser] = useRecoilState(userState);
  const [darkTheme, setDarkTheme] = useRecoilState(themeState);

  useEffect(() => {
    if (user) return;

    fetchAPI<MeUser>("/users/@me").then(
      (res) => setUser(res),
      (err) => {
        console.log("fetching /users/@me", err);
        if (
          (err as APIError).code == ErrorCode.InvalidToken ||
          (err as APIError).code == ErrorCode.Forbidden
        ) {
          localStorage.removeItem("pronouns-token");
        }
      }
    );
  }, [user, setUser]);

  const [showMenu, setShowMenu] = useState(false);

  useEffect(() => {
    setDarkTheme(
      localStorage.theme === "dark" ||
        (!("theme" in localStorage) &&
          window.matchMedia("(prefers-color-scheme: dark)").matches)
    );

    if (darkTheme) {
      document.documentElement.classList.add("dark");
    } else {
      document.documentElement.classList.remove("dark");
    }
  }, [darkTheme]);

  const storeTheme = (useDarkTheme: boolean | null) => {
    if (useDarkTheme === null) {
      localStorage.removeItem("theme");
    } else {
      localStorage.setItem("theme", useDarkTheme ? "dark" : "light");
    }
  };

  const nav = user ? (
    <>
      <NavItem href={`/u/${user.name}`}>@{user.name}</NavItem>
      <NavItem href="/settings">Settings</NavItem>
      <NavItem href="/logout">Log out</NavItem>
    </>
  ) : (
    <>
      <NavItem href="/login">Log in</NavItem>
    </>
  );

  return (
    <>
      <div className="bg-white/75 dark:bg-slate-800/75 w-full backdrop-blur border-slate-200 dark:border-slate-700 border-b">
        <div className="max-w-8xl mx-auto">
          <div className="py-4 mx-4">
            <div className="flex items-center">
              <Link href="/" passHref>
                <Logo />
              </Link>
              <div className="ml-auto flex items-center">
                <nav className="hidden lg:flex">
                  <ul className="flex space-x-4 font-bold">{nav}</ul>
                </nav>
                <div className="flex border-l border-slate-200 ml-4 pl-4 lg:ml-6 lg:pl-6 lg:mr-2 dark:border-slate-700 space-x-2 lg:space-x-4">
                  <div
                    onClick={() => {
                      setDarkTheme(!darkTheme);
                      storeTheme(!darkTheme);
                    }}
                    title={
                      darkTheme ? "Switch to light mode" : "Switch to dark mode"
                    }
                    className="cursor-pointer"
                  >
                    {darkTheme ? (
                      <Sun className="hover:text-sky-400" size={24} />
                    ) : (
                      <MoonStars size={24} className="hover:text-sky-500" />
                    )}
                  </div>
                  <div
                    onClick={() => setShowMenu(!showMenu)}
                    title="Show menu"
                    className="cursor-pointer flex lg:hidden"
                  >
                    <List
                      className="dark:hover:text-sky-400 hover:text-sky-500"
                      size={24}
                    />
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
      <nav
        className={`lg:hidden p-4 border-slate-200 dark:border-slate-700 border-b ${
          showMenu ? "flex" : "hidden"
        }`}
      >
        <ul className="flex flex-col space-y-4 font-bold">{nav}</ul>
      </nav>
    </>
  );
}
