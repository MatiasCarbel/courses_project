import { NextRequest, NextResponse } from "next/server";
import { jwtDecode } from "jwt-decode";
import { cookies } from "next/headers";

export async function GET(request: NextRequest) {
  const baseUrl = process.env.NEXT_PUBLIC_COURSES_API_URL ?? "";
  const cookieStore = cookies();
  const token = cookieStore.get("token");

  if (!token?.value) {
    return NextResponse.json({ message: "Not authenticated" }, { status: 401 });
  }

  try {
    const decoded = jwtDecode(token.value) as any;
    if (!decoded?.admin) {
      return NextResponse.json({ message: "Not authorized" }, { status: 403 });
    }

    const containersReq = await fetch(`${baseUrl}/containers`, {
      headers: {
        Cookie: `token=${token.value}`,
      },
      credentials: "include",
    });

    if (!containersReq.ok) {
      throw new Error(
        `Failed to fetch containers: ${containersReq.statusText}`
      );
    }

    const containers = await containersReq.json();
    return NextResponse.json(containers);
  } catch (error: any) {
    console.error("Container fetch error:", error);
    return NextResponse.json(
      { message: "Error fetching containers", error: error.message },
      { status: 500 }
    );
  }
}
