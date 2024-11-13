"use client";
import { CardTitle, CardDescription, CardHeader, CardContent, CardFooter, Card } from "@/components/ui/card"
import { Label } from "@/components/ui/label"
import { Input } from "@/components/ui/input"
import { Button } from "@/components/ui/button"
import Image from "next/image"
import S from "./Register.module.css";
import Link from "next/link"
import { useState } from "react"
import { useRouter } from "next/navigation";
import routes from "@/lib/routes";

export default function Component() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [passwordVerification, setPasswordVerification] = useState("");
  const [username, setUsername] = useState("");
  const [error, setError] = useState("");
  const router = useRouter();

  const handleSubmit = async () => {
    setError("");

    // check that email is well formatted.
    if (!email.includes("@") || !email.includes(".")) {
      setError("Invalid email format");
      return;
    }

    // check that password is at least 8 characters long.
    if (password.length < 8) {
      setError("Password must be at least 8 characters long");
      return;
    }

    // Check that both passwords match
    if (password !== passwordVerification) {
      setError("Passwords do not match");
      return;
    }

    // Check username is not empty
    if (username.trim() === "") {
      setError("Username cannot be empty");
      return;
    }

    try {
      const response = await fetch("/api/auth/register", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          username,
          email,
          password,
        }),
      });

      const data = await response.json();

      if (!response.ok) {
        setError(data.message || "Registration failed");
        return;
      }

      // If successful, redirect to login
      router.push(routes.login);
    } catch (err) {
      setError("An error occurred during registration");
    }
  };

  return (
    <div className="flex items-center justify-center">
      <Card className="w-full max-w-md">
        <div className="flex justify-center py-6">
          <div className={`${S.titleContainer} ml-2`}>
            <h1 className={`text-2xl font-bold ${S.title}`}>UCC</h1>
            <Image src="/triangulito.png" alt="UCCedemy" width="15" height={20} className={S.triangulito} />
          </div>
        </div>
        <CardHeader>
          <CardTitle className="text-[rgb(159,51,233)]">Register</CardTitle>
          <CardDescription>Enter your details to create a new account.</CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          {error && (
            <div className="p-3 text-sm text-red-500 bg-red-50 rounded-md">
              {error}
            </div>
          )}
          <div className="space-y-2">
            <Label htmlFor="username">Username</Label>
            <Input
              id="username"
              placeholder="johndoe"
              required
              value={username}
              onChange={e => setUsername(e.target.value)}
            />
          </div>
          <div className="space-y-2">
            <Label htmlFor="email">Email</Label>
            <Input
              id="email"
              placeholder="john@example.com"
              required
              value={email}
              onChange={e => setEmail(e.target.value)}
              type="email"
            />
          </div>
          <div className="space-y-2">
            <Label htmlFor="password">Password</Label>
            <Input
              id="password"
              required
              value={password}
              onChange={e => setPassword(e.target.value)}
              type="password"
            />
          </div>
          <div className="space-y-2">
            <Label htmlFor="confirm-password">Confirm Password</Label>
            <Input
              id="confirm-password"
              required
              value={passwordVerification}
              onChange={e => setPasswordVerification(e.target.value)}
              type="password"
            />
          </div>
        </CardContent>
        <CardFooter className="grid gap-2">
          <Button
            className="w-full bg-[rgb(159,51,233)] text-white hover:bg-[rgb(159,51,233)]/90"
            onClick={handleSubmit}
            type="button"
          >
            Register
          </Button>
          <Link href="/login">
            <Button
              className="w-full border-[rgb(159,51,233)] text-[rgb(159,51,233)] hover:text-[rgb(159,51,233)] hover:bg-[rgb(159,51,233)]/10"
              variant="outline"
            >
              Login
            </Button>
          </Link>
        </CardFooter>
      </Card>
    </div>
  );
}