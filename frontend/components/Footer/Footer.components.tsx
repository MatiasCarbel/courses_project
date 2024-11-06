"use client";

import Link from "next/link";
import S from "./Footer.module.css";
import Image from "next/image";

export default function Footer() {
  return (
    <footer className={`bg-[var(--gray)] ${S.Footer} flex items-center justify-between p-4 px-12 h-28`}>
      <Link href="/home" className={`${S.titleContainer}`}>
        <h2 className={`text-xl font-bold ${S.title}`}>UCCedemy</h2>
        <Image src="/triangulito.png" alt="UCCedemy" width="15" height={20} className={S.triangulito} />
      </Link>
      
      <p className="text-xs text-white">© 2024 UCCedemy, Inc.</p>
    </footer>
  );
}