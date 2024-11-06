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
  const [Email, setEmail] = useState("");
  const [Password, setPassword] = useState("");
  const [PasswordVerification, setPasswordVerification] = useState("");
  const [Username, setUsername] = useState("");
  const [FirstName, setFirstName] = useState("");
  const [LastName, setLastName] = useState("");
  const [UserType, setUserType] = useState(false);
  const router = useRouter();

  const handleSubmit = async () => {
    // check that Email is well formated.
    if (!Email.includes("@") && !Email.includes(".")) {
      console.error("Invalid Email");
      return;
    }

    // check that Password is at least 8 characters long.
    if (Password.length < 8) {
      console.error("Password must be at least 8 characters long");
      return;
    }

    // Check that both passwords match
    if (Password !== PasswordVerification) {
      console.error("Passwords do not match");
      return;
    }

    // Check name, lastname and username are not empty.
    if (FirstName === "" || LastName === "" || Username === "") {
      console.error("Name, lastname and username cannot be empty");
      return;
    }

    // Handle login logic here
    const response = await fetch("/api/auth/register", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ Email, PasswordHash: Password, Username, FirstName, LastName, UserType }),
    });

    const responseJson = await response.json();
    console.log(responseJson);

    if (response.ok) {
      router.push(routes.login);
    } else {
      console.error("Invalid Register");
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
          <div className="flex items-center justify-between">
            <div className="space-y-2">
              <Label htmlFor="name">
                Name
              </Label>
              <Input id="name" placeholder="John" required value={FirstName} onChange={e => setFirstName(e.target.value)} type="text" />
            </div>

            <div className="space-y-2">
              <Label htmlFor="name">
                Lastname
              </Label>
              <Input id="name" placeholder="Doe" required value={LastName} onChange={e => setLastName(e.target.value)} type="text" />
            </div>
          </div>
          <div className="space-y-2">
            <Label htmlFor="name">
              Username
            </Label>
            <Input id="name" placeholder="Johny" required value={Username} onChange={e => setUsername(e.target.value)} type="text" />
          </div>
          <div className="space-y-2">
            <Label htmlFor="email">
              Email
            </Label>
            <Input id="email" placeholder="m@example.com" required value={Email} onChange={e => setEmail(e.target.value)} type="email" />
          </div>
          <div className="space-y-2">
            <div className="flex items-center justify-between">
              <Label htmlFor="password">
                Password
              </Label>
            </div>
            <Input id="password" required value={Password} onChange={e => setPassword(e.target.value)} type="password" />
          </div>
          <div className="space-y-2">
            <div className="flex items-center justify-between">
              <Label htmlFor="confirm-password">
                Confirm Password
              </Label>
            </div>
            <Input id="confirm-password" required value={PasswordVerification} onChange={e => setPasswordVerification(e.target.value)} type="password" />
          </div>
        </CardContent>
        <CardFooter className="grid gap-2">
          <Button className="w-full bg-[rgb(159,51,233)] text-white hover:bg-[rgb(159,51,233)]/90" type="submit" onClick={handleSubmit}>
            Register
          </Button>
          <Link href="/login">
            <Button
              className="w-full border-[rgb(159,51,233)] text-[rgb(159,51,233)] hover:bg-[rgb(159,51,233)]/10 hover:text-[rgb(159,51,233)]"
              variant="outline"
            >
              Sign In
            </Button>
          </Link>
        </CardFooter>
      </Card>
    </div>
  )
}