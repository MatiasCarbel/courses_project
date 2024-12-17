import { NextResponse } from "next/server";
import { cookies } from "next/headers";

export async function GET() {
  const response = NextResponse.json(
    { message: "Logged out successfully" },
    { status: 200 }
  );

  // Remove the auth cookie
  response.cookies.set({
    name: "token",
    value: "",
    expires: new Date(0),
    path: "/",
  });

  return response;
}
