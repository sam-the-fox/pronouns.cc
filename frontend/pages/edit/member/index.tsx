import { useRouter } from "next/router";
import { useEffect } from "react";
import Loading from "../../../components/Loading";

export default function Redirect() {
    const router = useRouter();
    useEffect(() => {
        router.push("/")
    }, [])

    return <Loading />;
}