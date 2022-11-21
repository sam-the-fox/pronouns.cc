import Toastify from "toastify-js";
import "toastify-js/src/toastify.css";

export default function toast(options: { text: string; background?: string }) {
  let background: string;
  switch (options.background) {
    case "error":
      background = "#A1081F";
      break;
    case "success":
      background = "#1D611A";
      break;
    default:
      background = "#4F5859";
      break;
  }

  Toastify({
    text: options.text,
    gravity: "top",
    position: "left",
    duration: -1,
    close: true,
    style: {
      background: background,
      color: "#FFFFFF",
    },
  }).showToast();
}
