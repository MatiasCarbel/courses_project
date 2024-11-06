"use client";

import { Label } from "@/components/ui/label"
import { Input } from "@/components/ui/input"
import Link from "next/link"
import S from "./Login.module.css";
import { CardContent, CardFooter, Card, CardHeader, CardTitle, CardDescription } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import Image from "next/image"
import { useState } from "react";
import routes from "@/lib/routes";
import { useRouter } from "next/navigation";
import { useUser } from "@/hooks/useUser";

export default function Login() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const router = useRouter();
  const { refreshUser } = useUser();

  const handleSubmit = async () => {
    // check that email is well formated.
    if (!email.includes("@") && !email.includes(".")) {
      console.error("Invalid email");
      return;
    }

    // check that password is at least 8 characters long.
    if (password.length < 8) {
      console.error("Password must be at least 8 characters long");
      return;
    }

    // Handle login logic here
    const response = await fetch("/api/auth/login", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ email, password }),
    });

    const responseJson = await response.json();
    console.log(responseJson);

    if (response.ok) {
      await refreshUser();
      router.push(routes.home);
    } else {
      // Invalid credentials
      console.error("Invalid credentials");
    }
  };

  return (
    <div className="flex items-center justify-center">
      <Card className="w-full max-w-md ">
        <div className="flex justify-center py-6">
          <div className={`${S.titleContainer} ml-2`}>
            <h1 className={`text-2xl font-bold ${S.title}`}>UCC</h1>
            <Image src="/triangulito.png" alt="UCCedemy" width="15" height={20} className={S.triangulito} />
          </div>
        </div>
        <CardHeader>
          <CardTitle className="text-[rgb(159,51,233)]">Login</CardTitle>
          <CardDescription>Enter your email and password to sign in to your account.</CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="email">
              Email
            </Label>
            <Input id="email" placeholder="m@example.com" value={email} onChange={e => setEmail(e.target.value)} required type="email" />
          </div>
          <div className="space-y-2">
            <div className="flex items-center justify-between">
              <Label htmlFor="password">
                Password
              </Label>
              <Link
                className="text-sm font-medium text-[rgb(159,51,233)] hover:underline dark:text-[rgb(159,51,233)]"
                href="#"
              >
                Forgot password?
              </Link>
            </div>
            <Input id="password" required value={password} onChange={e => setPassword(e.target.value)} type="password" />
          </div>
        </CardContent>
        <CardFooter className="grid gap-2">
          <Button className="w-full bg-[rgb(159,51,233)] text-white hover:bg-[rgb(159,51,233)]/90" type="submit" onClick={handleSubmit}>
            Sign In
          </Button>
          <Link href="/register">
            <Button
              className="w-full border-[rgb(159,51,233)] text-[rgb(159,51,233)] hover:text-[rgb(159,51,233)] hover:bg-[rgb(159,51,233)]/10"
              variant="outline"
            >
              Register
            </Button>
          </Link>
        </CardFooter>
      </Card>
    </div>
  )
}
