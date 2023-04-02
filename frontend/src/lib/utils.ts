import MarkdownIt from "markdown-it";
import sanitize from "sanitize-html";

const md = new MarkdownIt({
  html: false,
  breaks: true,
}).disable(["heading", "link", "table"]);

export function renderMarkdown(src: string | null) {
  return src ? sanitize(md.render(src)) : null;
}

export const charCount = (str: string) => [...str].length;
