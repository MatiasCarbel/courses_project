import Link from "next/link";
import S from "./Button.module.css";
import { MouseEventHandler } from "react";

export default function Button({ children, variant, path, handleClick }: Readonly<{ children: React.ReactNode, variant: 'primary' | 'secondary', path?: string, handleClick?: MouseEventHandler<HTMLAnchorElement> | undefined }>) {
  return (
      <Link href={path ?? ''} className={`${S.button} ${S[variant]} ml-8`} onClick={handleClick}>
        <p className="text-sm font-bold">{children}</p>
      </Link>
  );
}