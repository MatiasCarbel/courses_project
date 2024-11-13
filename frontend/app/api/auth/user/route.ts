import { NextResponse } from "next/server";
import { cookies } from "next/headers";
import { jwtDecode } from "jwt-decode";

export async function GET(request: Request) {
  const cookiesStore = cookies();
  const cookie = cookiesStore.get("auth");

  const cookieValue = cookie?.value;
  if (!cookieValue) {
    return NextResponse.json(
      { message: "Not authenticated.", shouldLogin: true },
      { status: 200 }
    );
  }

  const decoded = jwtDecode(cookieValue) as any;
  const userId = decoded?.user_id;

  if (!userId) {
    return NextResponse.json(
      { message: "Invalid token format", shouldLogin: true },
      { status: 200 }
    );
  }

  const baseUrl =
    process.env.NEXT_PUBLIC_USERS_API_URL ?? "http://users-api:8001";

  console.log("baseUrl: ", baseUrl);
  console.log("userId: ", userId);
  console.log("cookieValue: ", cookieValue);

  try {
    const userReq = await fetch(`${baseUrl}/users/${userId}`, {
      headers: {
        Authorization: `Bearer ${cookieValue}`,
      },
    });

    if (!userReq.ok) {
      return NextResponse.json(
        { message: "Failed to fetch user", shouldLogin: true },
        { status: 200 }
      );
    }

    const userJson = await userReq.json();
    return NextResponse.json(
      { message: "OK", user: userJson },
      { status: 200 }
    );
  } catch (error) {
    console.error("Error fetching user:", error);
    return NextResponse.json(
      { message: "Error fetching user", shouldLogin: true },
      { status: 200 }
    );
  }
}
