import { HTMLAttributes } from "react";

export interface Props extends HTMLAttributes<Props> {
  urls: string[];
  alt: string;
}

export default function FallbackImage({ urls, alt, className }: Props) {
  const fallbackUrl = urls.pop()!;
  urls.push(fallbackUrl);

  return (
    <picture className={className}>
      {urls.length !== 0 &&
        urls.map((url, key) => {
          let contentType: string;
          if (url.endsWith(".webp")) {
            contentType = "image/webp";
          } else if (url.endsWith(".jpg") || url.endsWith(".jpeg")) {
            contentType = "image/jpeg";
          } else if (url.endsWith(".png")) {
            contentType = "image/png";
          } else {
            contentType = "application/octet-stream";
          }
          return <source key={key} srcSet={url} type={contentType} />;
        })}
      <img src={fallbackUrl} alt={alt} className={className} />
    </picture>
  );
}
