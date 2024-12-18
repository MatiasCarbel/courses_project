import { NextResponse } from "next/server";

export async function POST(request: Request) {
  const data = await request.json();
  const { email, password } = data;

  const baseUrl =
    process.env.NEXT_PUBLIC_USERS_API_URL ?? "http://users-api:8001";

  try {
    const loginReq = await fetch(`${baseUrl}/user/login`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ email, password }),
      credentials: "include",
    });

    if (!loginReq.ok) {
      const errorText = await loginReq.text();
      console.error("Login failed:", errorText);
      return NextResponse.json(
        { message: "Invalid credentials" },
        { status: 401 }
      );
    }

    const loginJson = await loginReq.json();

    // Create response with cookie
    const response = NextResponse.json(
      { message: "Logged In.", user: loginJson },
      { status: 200 }
    );

    // Set the cookie from the API response
    response.cookies.set({
      name: "token",
      value: loginJson.token,
      httpOnly: false,
      secure: process.env.NODE_ENV === "production",
      sameSite: "lax",
      path: "/",
    });

    return response;
  } catch (error) {
    console.error("Login error:", error);
    return NextResponse.json(
      { message: "An error occurred during login" },
      { status: 500 }
    );
  }
}
