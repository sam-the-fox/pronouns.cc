import ReactMarkdown from "react-markdown";

// this is a temporary home page, which is why the markdown content is embedded
const md = `This will (one day) be a site to create pronoun cards for yourself,
similarly to [Pronouny](https://pronouny.xyz/) and [Pronouns.page](https://en.pronouns.page/).

You'll be able to create multiple profiles that are linked together,
useful for plurality ([what?](https://morethanone.info/)) and kin, or even just for fun!

For now though, there's just this landing page <3  
(And no, the "Log in" button doesn't do anything either.)

Check out the (work in progress) source code on [GitLab](https://gitlab.com/1f320/pronouns)!`;

function Home() {
  return (
    <div className="prose prose-slate dark:prose-invert">
      <ReactMarkdown>{md}</ReactMarkdown>
    </div>
  );
}

export default Home;
