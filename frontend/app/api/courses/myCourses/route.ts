import { formatCookies, getCookieValue } from "@/lib/api.utils";
import { jwtDecode } from "jwt-decode";
import { NextRequest, NextResponse } from "next/server";

export async function GET(request: NextRequest) {
  const baseUrl = process.env.NEXT_PUBLIC_BASE_API_URL ?? "";

  const cookie = formatCookies();
  const cookieValue = getCookieValue();

  const decoded = jwtDecode(cookieValue ?? "") as any;
  const userId = decoded?.user_id;

  const url = `${baseUrl}/user/courses/${userId}`;

  const coursesReq = await fetch(url, {
    headers: {
      cookie: cookie,
    },
  });

  const coursesJson = await coursesReq.json();
  const courses = coursesJson.results;

  return NextResponse.json(
    { message: "Courses fetched", courses },
    { status: 200 }
  );
}
