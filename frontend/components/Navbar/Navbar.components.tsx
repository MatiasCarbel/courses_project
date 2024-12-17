"use client";

import Link from "next/link";
import Button from "../Button/Button.components";
import S from "./Navbar.module.css";
import Image from "next/image";
import { useUser } from "@/hooks/useUser";

export default function Navbar() {
  const { user, isAdmin, isAuthed, logout, isLoading } = useUser();

  if (isLoading) {
    return (
      <header className={`${S.navbar} flex items-center justify-between p-4`}>
        <Link href="/home" className={`${S.titleContainer} ml-2`}>
          <h1 className={`text-2xl font-bold ${S.title}`}>UCCedemy</h1>
          <Image src="/triangulito.png" alt="UCCedemy" width="15" height={20} className={S.triangulito} />
        </Link>
      </header>
    );
  }

  return (
    <header className={` ${S.navbar} flex items-center justify-between p-4`}>
      <Link href="/home" className={`${S.titleContainer} ml-2`}>
        <h1 className={`text-2xl font-bold ${S.title}`}>UCCedemy</h1>
        <Image src="/triangulito.png" alt="UCCedemy" width="15" height={20} className={S.triangulito} />
      </Link>
      <nav className="flex items-center">
        {isAdmin && (
          <Link href="/upload" className="ml-8">
            <p className="text-sm font-bold">Create Course</p>
          </Link>
        )}

        <Link href="/home" className="ml-8">
          <p className="text-sm font-bold">Courses</p>
        </Link>

        {isAuthed && (
          <>
            <Link href="/myCourses" className="ml-8">
              <p className="text-sm font-bold">My Courses</p>
            </Link>

            <Button handleClick={() => logout()} variant="secondary">{user?.username}</Button>
          </>
        )}

        {!isAuthed && (
          <>
            <Button path="/login" variant="primary">Login</Button>

            <Button path="/register" variant="secondary">Sign Up</Button>
          </>
        )}

        {isAdmin && (
          <Link href="/admin/services" className="ml-8">
            <p className="text-sm font-bold">Services Dashboard</p>
          </Link>
        )}
      </nav>
    </header>
  );
}