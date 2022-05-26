import { ThreeDots } from "react-bootstrap-icons";

export default function Loading() {
  return (
    <div className="flex flex-col pt-32 items-center">
      <ThreeDots size={64} className="animate-bounce" aria-hidden="true" />
      <span className="font-bold text-xl">Loading...</span>
    </div>
  );
}
