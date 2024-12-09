import { NextResponse } from "next/server";

export async function POST(request: Request) {
  const data = await request.json();
  const { username, email, password } = data;

  const baseUrl =
    process.env.NEXT_PUBLIC_USERS_API_URL ?? "http://users-api:8001";

  try {
    const usersReq = await fetch(`${baseUrl}/users`, {
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

    if (!usersReq.ok) {
      const errorData = await usersReq.json();
      return NextResponse.json(
        { message: errorData.error || "Registration failed" },
        { status: usersReq.status }
      );
    }

    return NextResponse.json(
      { message: "User created successfully" },
      { status: 200 }
    );
  } catch (error) {
    console.error("Registration error:", error);
    return NextResponse.json(
      { message: "An error occurred during registration" },
      { status: 500 }
    );
  }
}
