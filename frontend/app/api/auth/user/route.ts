import { NextResponse } from "next/server";
import { cookies } from "next/headers";
import { jwtDecode } from "jwt-decode";

export async function GET(request: Request) {
  const cookiesStore = cookies();
  const cookie = cookiesStore.get("auth");

  const cookieValue = cookie?.value;
  if (!cookieValue)
    return NextResponse.json(
      { message: "Not authenticated.", shouldLogin: true },
      { status: 200 }
    );

  const decoded = jwtDecode(cookieValue ?? "") as any;
  const userId = decoded?.id;

  const baseUrl = process.env.NEXT_PUBLIC_BASE_API_URL ?? "";
  const userReq = await fetch(`${baseUrl}/user/${userId}`, {
    headers: {
      cookie: `auth=${cookieValue}`,
    },
  });

  const userJson = await userReq.json();
  if (userJson.error)
    return NextResponse.json(
      { message: userJson.error, shouldLogin: true },
      { status: 200 }
    );

  return NextResponse.json({ message: "OK", user: userJson }, { status: 200 });
}
