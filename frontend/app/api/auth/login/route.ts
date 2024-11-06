import { NextResponse } from "next/server";
import { cookies } from "next/headers";

export async function POST(request: Request) {
  // Handle login logic here
  const data = await request.json();
  const { email, password } = data;

  const baseUrl = process.env.NEXT_PUBLIC_BASE_API_URL ?? "";

  const usersReq = await fetch(`${baseUrl}/user/login`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ email, password }),
  });

  const userJson = await usersReq.json();
  if (userJson.error)
    return NextResponse.json({ message: userJson.error }, { status: 401 });

  const cookieStore = cookies();
  cookieStore.set("auth", userJson?.token ?? "");

  return NextResponse.json(
    { message: "Logged In.", user: userJson },
    { status: 200 }
  );
}
