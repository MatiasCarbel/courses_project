import { NextRequest, NextResponse } from "next/server";
import { formatCookies, getCookieValue } from "@/lib/api.utils";
import { jwtDecode } from "jwt-decode";

export async function GET() {
  const baseUrl = process.env.NEXT_PUBLIC_COURSES_API_URL ?? "";
  const cookie = formatCookies();
  const cookieValue = getCookieValue();

  if (!cookieValue) {
    return NextResponse.json({ error: "Not authenticated" }, { status: 401 });
  }

  const decoded = jwtDecode(cookieValue) as any;
  if (!decoded?.admin) {
    return NextResponse.json({ error: "Not authorized" }, { status: 403 });
  }

  try {
    const response = await fetch(`${baseUrl}/api/services`, {
      headers: {
        Cookie: cookie,
        Authorization: `Bearer ${cookieValue}`,
      },
      credentials: "include",
    });
    const data = await response.json();

    if (!response.ok) {
      throw new Error(data.error || "Failed to fetch services");
    }

    return NextResponse.json(data);
  } catch (error: any) {
    return NextResponse.json(
      { error: error.message || "Failed to fetch services" },
      { status: 500 }
    );
  }
}

export async function POST(request: NextRequest) {
  const baseUrl = process.env.NEXT_PUBLIC_COURSES_API_URL ?? "";
  const cookie = formatCookies();
  const cookieValue = getCookieValue();

  if (!cookieValue) {
    return NextResponse.json({ error: "Not authenticated" }, { status: 401 });
  }

  const decoded = jwtDecode(cookieValue) as any;
  if (!decoded?.admin) {
    return NextResponse.json({ error: "Not authorized" }, { status: 403 });
  }

  try {
    const body = await request.json();
    const response = await fetch(`${baseUrl}/api/services`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Cookie: cookie,
        Authorization: `Bearer ${cookieValue}`,
      },
      credentials: "include",
      body: JSON.stringify(body),
    });

    const data = await response.json();

    if (!response.ok) {
      throw new Error(data.error || "Failed to add instance");
    }

    return NextResponse.json(data);
  } catch (error: any) {
    return NextResponse.json(
      { error: error.message || "Failed to add instance" },
      { status: 500 }
    );
  }
}
