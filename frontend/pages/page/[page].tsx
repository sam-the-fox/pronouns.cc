import { GetStaticProps } from "next";
import { readdirSync } from "fs";
import ReactMarkdown from "react-markdown";
import { readFile } from "fs/promises";
import { join } from "path";
import Head from "next/head";

export async function getStaticPaths() {
  const names = readdirSync("./static_pages").filter((name) =>
    name.endsWith(".md")
  );

  const paths = names.map((name) => ({
    params: { page: name.slice(0, -3) },
  }));

  return {
    paths: paths,
    fallback: false,
  };
}

export const getStaticProps: GetStaticProps<{ text: string }> = async ({
  params,
}) => {
  const text = await readFile(join("./static_pages", params!.page + ".md"));
  return { props: { text: text.toString("utf8") } };
};

export default function Page(props: { text: string }) {
  const title = props.text.split("\n")[0].slice(2);

  return (
    <>
      <Head>
        <title key="title">{`${title} - pronouns.cc`}</title>
      </Head>
      <div className="prose prose-slate dark:prose-invert">
        <ReactMarkdown>{props.text}</ReactMarkdown>
      </div>
    </>
  );
}
