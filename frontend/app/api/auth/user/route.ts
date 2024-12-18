import { NextResponse } from "next/server";
import { cookies } from "next/headers";
import { jwtDecode } from "jwt-decode";

export async function GET(request: Request) {
  const cookiesStore = cookies();
  const cookie = cookiesStore.get("token");

  const cookieValue = cookie?.value;
  if (!cookieValue) {
    return NextResponse.json(
      { message: "Not authenticated", user: null },
      { status: 200 }
    );
  }

  const decoded = jwtDecode(cookieValue);

  // Return the decoded token data as the user object
  return NextResponse.json({ message: "OK", user: decoded }, { status: 200 });
}
